package controller

import (
	"log/slog"
	"sync"

	"github.com/e-berger/sheepdog-runner/internal/metrics"
	"github.com/e-berger/sheepdog-runner/internal/probes"
)

type Controller struct {
	pushGateway *metrics.Push
}

func NewController(pushGateway string) *Controller {
	var p *metrics.Push
	if pushGateway != "" {
		p = metrics.NewPush(pushGateway)
	}
	return &Controller{
		pushGateway: p,
	}
}

func (c *Controller) Run(probesDatas []probes.IProbe) (int, int, error) {
	monitorErr := 0
	wg := new(sync.WaitGroup)
	for _, probe := range probesDatas {
		wg.Add(1)
		go func(probe probes.IProbe) {
			slog.Info("Launching monitoring", "probe", probe.String())
			defer wg.Done()
			result, err := probe.Launch()
			if err != nil {
				monitorErr++
				slog.Error("Error launching monitoring", "error", err)
			}
			slog.Info("Metrics monitoring", "probe", result.String())
			if c.pushGateway != nil {
				err = c.SendMetrics(result)
				if err != nil {
					monitorErr++
					slog.Error("Error pushing monitoring", "error", err)
				}
			} else {
				slog.Info("No push gateway defined")
			}
		}(probe)
	}
	wg.Wait()
	slog.Info("Monitor", "nb error", monitorErr)
	return len(probesDatas), monitorErr, nil
}

func (c *Controller) SendMetrics(metrics metrics.IMetrics) error {
	return c.pushGateway.Send(metrics.GetId(), metrics.GetMetrics())
}
