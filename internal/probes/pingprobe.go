package probes

import (
	"fmt"

	domain "github.com/e-berger/sheepdog-domain/probes"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/e-berger/sheepdog-runner/internal/results"
)

type pingProbe struct {
	domain.Probe
	location types.Location
}

func (t pingProbe) String() string {
	return fmt.Sprintf("http probe %s", t.Probe.GetId())
}

func (t pingProbe) GetId() string {
	return t.Probe.GetId()
}

func (t pingProbe) IsInError() bool {
	return t.Probe.IsInError()
}

func NewPingProbe(probe domain.Probe, location types.Location) (IProbe, error) {
	return pingProbe{
		probe,
		location,
	}, nil
}

func (t pingProbe) GetHttpClient() HTTPClient {
	return nil
}

func (t pingProbe) Launch(client HTTPClient) results.IResults {
	result := results.NewResultsPingEmpty(t.GetId(), t.location)
	return result
}
