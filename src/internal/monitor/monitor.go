package monitor

import (
	"fmt"

	"github.com/e-berger/sheepdog-runner/src/internal/monitoring"
)

type IMonitor interface {
	GetMethod() monitoring.MethodType
	Launch() error
	String() string
}

type Monitor struct {
	id       int
	method   monitoring.MethodType
	location monitoring.LocationType
}

func (t *Monitor) GetMethod() monitoring.MethodType {
	return t.method
}

func (t *Monitor) Launch() error {
	return nil
}

func (t *Monitor) String() string {
	return fmt.Sprintf("monitor %d", t.id)
}

func GetMonitoring(m *InfosMonitor) (IMonitor, error) {
	switch {
	case m.Type == monitoring.HTTP:
		return newHttpTest(m)
	case m.Type == monitoring.PING:
		return newPingTest(m)
	default:
		return nil, fmt.Errorf("monitor type %d not found", m.Type)
	}
}
