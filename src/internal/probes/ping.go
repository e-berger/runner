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

func (t *pingProbe) Launch() (metrics.IMetrics, error) {
	// msg := &icmp.Message{
	// 	Type: ipv4.ICMPType,
	// 	Code: 0,
	// 	Body: &icmp.Echo{
	// 		ID:   os.Getpid() & 0xffff,
	// 		Seq:  0,
	// 		Data: []byte("hello"),
	// 	},
	// }
	// wb, err := msg.Marshal(nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return nil, nil
}

func (t *pingProbe) String() string {
	return fmt.Sprintf("ping probe %s", t.Id)
}

func (t *pingProbe) GetResult() metrics.IMetrics {
	return t.result
}
