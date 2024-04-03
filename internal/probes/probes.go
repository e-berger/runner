package probes

import (
	"encoding/json"
	"fmt"

	domain "github.com/e-berger/sheepdog-domain/probes"
	"github.com/e-berger/sheepdog-runner/internal/metrics"
)

const (
	DEFAULT_VALID = 1
)

type IProbe interface {
	Launch() (metrics.IMetrics, error)
	GetId() string
	IsInError() bool
	String() string
}

type ProbeJSON struct {
	Id   string          `json:"id"`
	Info json.RawMessage `json:"info"`
	Type int             `json:"type"`
}

type Probes struct {
	Location domain.Location `json:"location"`
	Probes   []IProbe        `json:"items"`
	Mode     domain.Mode     `json:"mode"`
}

func NewProbeFromJSON(probeJSON ProbeJSON, location string) (IProbe, error) {

	loc, err := domain.ParseLocation(location)
	if err != nil {
		return nil, err
	}
	locs, err := domain.ParseLocations([]string{location})
	if err != nil {
		return nil, err
	}

	switch domain.ProbeType(probeJSON.Type) {
	case domain.HttpProbeType:
		httpInfo, err := domain.NewHttpProbeInfoFromJson([]byte(probeJSON.Info))
		if err != nil {
			return nil, err
		}
		probe, err := domain.NewProbeHttp(probeJSON.Id, domain.DefaultInterval, locs, false, httpInfo)
		if err != nil {
			return nil, err
		}
		return NewHttpProbe(probe, loc)
	case domain.PingProbeType:
		pingInfo, err := domain.NewPingProbeInfoFromJson([]byte(probeJSON.Info))
		if err != nil {
			return nil, err
		}
		probe, err := domain.NewProbePing(probeJSON.Id, domain.DefaultInterval, locs, false, pingInfo)
		if err != nil {
			return nil, err
		}
		return NewPingProbe(probe, loc)
	default:
		return nil, fmt.Errorf("unknown probe type %d", probeJSON.Type)
	}
}
