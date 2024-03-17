package probes

import (
	"fmt"

	"github.com/e-berger/sheepdog-runner/internal/metrics"
)

type tcpProbe struct {
	Probe
}

func NewTcpProbe(p *Probe) (IProbe, error) {
	return &tcpProbe{
		Probe: Probe{
			Id:       p.Id,
			Type:     PING,
			Location: p.Location,
			State:    p.State,
		},
	}, nil
}

func (t *tcpProbe) Launch() (metrics.IMetrics, error) {
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

func (t *tcpProbe) String() string {
	return fmt.Sprintf("tcp probe %s", t.Id)
}

func (t *tcpProbe) IsError() bool {
	return false
}
