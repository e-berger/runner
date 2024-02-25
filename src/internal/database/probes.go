package database

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/e-berger/sheepdog-runner/src/internal/domain"
	"github.com/e-berger/sheepdog-runner/src/internal/probes"
)

func (t *TursoDatabase) GetProbes(limit int, offset int) ([]domain.IProbe, error) {
	statements := map[string]interface{}{
		"statements": []map[string]interface{}{
			{
				"q": "SELECT id, type, data, location FROM probes ORDER BY id LIMIT :limit OFFSET :offset",
				"params": map[string]interface{}{
					":limit":  limit,
					":offset": offset,
				},
			},
		},
	}

	time_start := time.Now()
	responseBody, err := t.httpQuery(statements)
	if err != nil {
		return nil, err
	}
	slog.Debug("Response", "query time", time.Since(time_start).String())

	var resultData []Response
	err = json.Unmarshal(responseBody, &resultData)
	if err != nil {
		return nil, err
	}

	probes, _ := parseProbes(resultData)
	slog.Debug("Response", "get probes result", probes)
	return probes, nil
}

func parseProbes(result []Response) ([]domain.IProbe, error) {
	var probesArray []domain.IProbe
	var p domain.IProbe
	for _, result := range result {
		for l := range result.Results.Rows {
			probe, err := domain.NewProbe(result.Results.Columns, result.Results.Rows[l])
			if err != nil {
				return nil, err
			}
			p, err = probes.CreateProbeFromType(probe)
			if err != nil {
				return nil, err
			}
			probesArray = append(probesArray, p)
		}
	}
	return probesArray, nil
}
