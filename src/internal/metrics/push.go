package metrics

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Push struct {
	pushGateway string
}

func NewPush(pushGateway string) *Push {
	return &Push{
		pushGateway: pushGateway,
	}
}

func (p *Push) Send(job string, completionTime prometheus.Collector) error {
	pusher := push.New(p.pushGateway, job)
	err := pusher.Collector(completionTime).Push()
	if err != nil {
		slog.Error("Could not push completion time to Pushgateway:", "error", err)
		return err
	}
	return nil
}
