package probes

import (
	"fmt"

	"github.com/e-berger/sheepdog-runner/src/internal/domain"
	"github.com/e-berger/sheepdog-runner/src/internal/metrics"
)

type pingProbe struct {
	domain.Probe
	result *metrics.MetricsHttp
}

func NewPingProbe(p *domain.Probe) (domain.IProbe, error) {
	return &pingProbe{
		Probe: domain.Probe{
			Id:       p.Id,
			Type:     domain.PING,
			Location: p.Location,
		},
	}, nil
}

func (t *pingProbe) GetType() domain.ProbeType {
	return t.Type
}

func (t *pingProbe) Launch() error {
	return nil
}

func (t *pingProbe) Push(pushGateway *metrics.Push) error {
	return nil
}

func (t *pingProbe) String() string {
	return fmt.Sprintf("ping test %s", t.Id)
}

func (t *pingProbe) GetResult() metrics.IMetrics {
	return t.result
}
