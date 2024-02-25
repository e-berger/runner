package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"log/slog"
)

type Response struct {
	Results Results `json:"results"`
}

type Results struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

type TursoDatabase struct {
	Database  string
	AuthToken string
}

func NewTursoDatabase(database string, authToken string) *TursoDatabase {
	return &TursoDatabase{
		Database:  fmt.Sprintf("https://%s.turso.io", database),
		AuthToken: authToken,
	}
}

func (t *TursoDatabase) httpQuery(statements map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(statements)
	if err != nil {
		slog.Error("Error marshalling JSON query", "error", err, "statements", statements)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, t.Database, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AuthToken))

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
