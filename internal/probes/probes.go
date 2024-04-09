package probes

import (
	"encoding/json"
	"fmt"

	domain "github.com/e-berger/sheepdog-domain/probes"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/e-berger/sheepdog-runner/internal/results"
)

const (
	DEFAULT_VALID = 1
)

type IProbe interface {
	GetHttpClient() HTTPClient
	Launch(client HTTPClient) results.IResults
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
	Location types.Location `json:"location"`
	Probes   []IProbe       `json:"items"`
	Mode     types.Mode     `json:"mode"`
}

func NewProbeFromJSON(probeJSON ProbeJSON, location string) (IProbe, error) {

	loc, err := types.ParseLocation(location)
	if err != nil {
		return nil, err
	}
	locs, err := types.ParseLocations([]string{location})
	if err != nil {
		return nil, err
	}

	switch types.ProbeType(probeJSON.Type) {
	case types.HttpProbeType:
		httpInfo, err := domain.NewHttpProbeInfoFromJson([]byte(probeJSON.Info))
		if err != nil {
			return nil, err
		}
		probe, err := domain.NewProbeHttp(probeJSON.Id, types.DefaultInterval, locs, false, httpInfo)
		if err != nil {
			return nil, err
		}
		return NewHttpProbe(probe, loc)
	case types.PingProbeType:
		pingInfo, err := domain.NewPingProbeInfoFromJson([]byte(probeJSON.Info))
		if err != nil {
			return nil, err
		}
		probe, err := domain.NewProbePing(probeJSON.Id, types.DefaultInterval, locs, false, pingInfo)
		if err != nil {
			return nil, err
		}
		return NewPingProbe(probe, loc)
	default:
		return nil, fmt.Errorf("unknown probe type %d", probeJSON.Type)
	}
}
