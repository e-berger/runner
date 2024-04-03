package probes

import (
	"fmt"
	"time"

	domain "github.com/e-berger/sheepdog-domain/probes"
	"github.com/e-berger/sheepdog-runner/internal/metrics"
)

type pingProbe struct {
	domain.Probe
	location domain.Location
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

func NewPingProbe(probe domain.Probe, location domain.Location) (IProbe, error) {
	return pingProbe{
		probe,
		location,
	}, nil
}

func (t pingProbe) Launch() (metrics.IMetrics, error) {
	time_start := time.Now()
	result := metrics.NewResultHttpDetails(
		t.GetId(),
		int(t.location),
		time_start,
		time.Since(time_start).Milliseconds(),
		DEFAULT_VALID,
		t.GetHttpProbeInfo().Method,
		200)
	return result, nil
}
