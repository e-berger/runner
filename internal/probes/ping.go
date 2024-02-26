package probes

import (
	"fmt"

	"github.com/e-berger/sheepdog-runner/internal/metrics"
)

type pingProbe struct {
	Probe
	result *metrics.MetricsHttp
}

func NewPingProbe(p *Probe) (IProbe, error) {
	return &pingProbe{
		Probe: Probe{
			Id:       p.Id,
			Type:     PING,
			Location: p.Location,
		},
	}, nil
}

func (t *pingProbe) GetType() ProbeType {
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
