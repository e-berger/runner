package datas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"log/slog"
	"time"

	"github.com/e-berger/sheepdog-runner/src/internal/monitor"
)

type Response struct {
	Results Results `json:"results"`
}

type Results struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

func Fetch(limit int, offset int, database string, authToken string) ([]*monitor.InfosMonitor, error) {
	statements := map[string]interface{}{
		"statements": []map[string]interface{}{
			{
				"q": "SELECT id, type, data, location FROM probes LIMIT :limit OFFSET :offset",
				"params": map[string]interface{}{
					":limit":  limit,
					":offset": offset,
				},
			},
		},
	}
	jsonData, err := json.Marshal(statements)
	if err != nil {
		slog.Error("Error marshalling JSON")
		return nil, err
	}

	time_start := time.Now()
	responseBody, err := httpQueryDb(database, authToken, jsonData)
	if err != nil {
		return nil, err
	}
	slog.Debug("Response", "query time", time.Since(time_start).String())

	var resultData []Response
	err = json.Unmarshal(responseBody, &resultData)
	if err != nil {
		return nil, err
	}

	resultDataFlatten, _ := parseData(resultData)
	slog.Debug("Response", "query result", resultDataFlatten)
	return resultDataFlatten, nil
}

func parseData(resultData []Response) ([]*monitor.InfosMonitor, error) {
	var resultDataFlatten []*monitor.InfosMonitor
	for _, result := range resultData {
		for l := range result.Results.Rows {
			infosMonitor, err := monitor.NewInfosMonitor(result.Results.Columns, result.Results.Rows[l])
			if err != nil {
				return nil, err
			}
			resultDataFlatten = append(resultDataFlatten, infosMonitor)
		}
	}
	return resultDataFlatten, nil
}

func httpQueryDb(database string, authToken string, jsonData []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, database, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		slog.Error("Errored when sending request to the server")
		return nil, err
	}
	slog.Debug("Response", "status", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}
