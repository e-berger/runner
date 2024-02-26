package probes

import (
	"fmt"
)

func CreateProbeFromType(p *Probe) (IProbe, error) {
	switch {
	case p.Type == HTTP:
		return NewHttpProbe(p)
	case p.Type == PING:
		return NewPingProbe(p)
	case p.Type == TCP:
		return NewTcpProbe(p)
	default:
		return nil, fmt.Errorf("probe type %d not found", p.Type)
	}
}
