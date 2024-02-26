package domain

import "github.com/e-berger/sheepdog-runner/internal/metrics"

type IProbe interface {
	GetType() ProbeType
	Launch() (metrics.IMetrics, error)
	String() string
	GetResult() metrics.IMetrics
}

type Probe struct {
	Id       string    `json:"id"`
	Type     ProbeType `json:"type"`
	Location Location  `json:"location"`
	Data     string    `json:"data"`
}

func NewProbe(columns []string, row []interface{}) (*Probe, error) {
	probe := &Probe{}
	for i, col := range columns {
		switch col {
		case "id":
			probe.Id = row[i].(string)
		case "type":
			probe.Type = ProbeType(int(row[i].(float64)))
		case "location":
			probe.Location = Location(int(row[i].(float64)))
		case "data":
			probe.Data = row[i].(string)
		}
	}
	return probe, nil
}
