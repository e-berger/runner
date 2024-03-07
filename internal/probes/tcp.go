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
			Error:    p.Error,
		},
	}, nil
}

func (t *tcpProbe) GetType() ProbeType {
	return t.Type
}

func (t *tcpProbe) GetId() string {
	return t.Id
}

func (t *tcpProbe) GetError() bool {
	return t.Error
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
