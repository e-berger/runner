package monitor

import (
	"fmt"

	"github.com/e-berger/sheepdog-runner/src/internal/monitoring"
)

type pingTest struct {
	Monitor
}

func newPingTest(m *InfosMonitor) (IMonitor, error) {
	return &pingTest{
		Monitor: Monitor{
			id:       m.Id,
			method:   monitoring.PING,
			location: monitoring.LocationType(m.Location),
		},
	}, nil
}

func (t *pingTest) GetMethod() monitoring.MethodType {
	return t.method
}

func (t *pingTest) Launch() error {
	return nil
}

func (t *pingTest) String() string {
	return fmt.Sprintf("ping test %d", t.id)
}
