package probes

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/e-berger/sheepdog-runner/internal/metrics"
)

type IProbe interface {
	GetType() ProbeType
	GetId() string
	GetMode() Mode
	IsError() bool
	Launch() (metrics.IMetrics, error)
	String() string
}

// {
// 	"id": "2d2f35Ry74DX9F9piVm4FWUuz3b",
// 	"info": {
// 		"timeout": 10000000000,
// 		"method": "GET",
// 		"url": "https://observations-service-api.eu.finalcad.cloud/healthz/live",
// 		"expected_status_code": [
// 			200
// 		]
// 	},
// 	"type": 2
// }

type Probe struct {
	Id       string          `json:"id"`
	Type     ProbeType       `json:"type"`
	Location Location        `json:"location,omitempty"`
	Mode     Mode            `json:"mode,omitempty"`
	Data     json.RawMessage `json:"info"`
	State    State           `json:"state"`
}

func UnmarshalJSON(data []byte, location string, mode string) (*Probe, error) {
	probe := &Probe{}
	slog.Info("test", "data", string(data))
	err := json.Unmarshal(data, probe)
	if err != nil {
		return nil, err
	}
	probe.Mode, err = ParseMode(mode)
	if err != nil {
		slog.Error("Error parsing mode", "mode", mode)
		return nil, err
	}
	probe.Location, err = ParseLocation(location)
	if err != nil {
		slog.Error("Error parsing location", "location", location)
		return nil, err
	}
	return probe, nil
}

func (p *Probe) CreateProbeFromType() (IProbe, error) {
	switch {
	case p.Type == HTTP:
		return NewHttpProbe(p)
	case p.Type == TCP:
		return NewTcpProbe(p)
	}
	return nil, fmt.Errorf("probe type %d not found", p.Type)
}

func (p *Probe) IsError() bool {
	return p.State == ERROR
}

func (p *Probe) GetType() ProbeType {
	return p.Type
}

func (p *Probe) GetId() string {
	return p.Id
}

func (p *Probe) GetMode() Mode {
	return p.Mode
}
