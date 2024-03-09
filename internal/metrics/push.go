package metrics

import (
	"context"
	"log/slog"
	"time"

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

func (p *Push) Send(id string, collector prometheus.Collector) error {
	pusher := push.New(p.pushGateway, id)
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*3)
	defer cncl()
	err := pusher.Collector(collector).PushContext(ctx)
	if err != nil {
		slog.Error("Could not push completion time to Pushgateway:", "error", err)
		return err
	}
	return nil
}
