package probes

import (
	"fmt"

	"github.com/e-berger/sheepdog-runner/src/internal/domain"
)

func CreateProbeFromType(p *domain.Probe) (domain.IProbe, error) {
	switch {
	case p.Type == domain.HTTP:
		return NewHttpProbe(p)
	case p.Type == domain.PING:
		return NewPingProbe(p)
	default:
		return nil, fmt.Errorf("probe type %d not found", p.Type)
	}
}
