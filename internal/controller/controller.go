package controller

import (
	"log/slog"
	"sync"

	db "github.com/e-berger/sheepdog-runner/internal/database"
	"github.com/e-berger/sheepdog-runner/internal/metrics"
	"github.com/e-berger/sheepdog-runner/internal/probes"
)

type Controller struct {
	Database    *db.TursoDatabase
	pushGateway *metrics.Push
}

func NewController(database string, authToken string, pushGateway string) *Controller {
	var p *metrics.Push
	if pushGateway != "" {
		p = metrics.NewPush(pushGateway)
	}
	return &Controller{
		Database:    db.NewTursoDatabase(database, authToken),
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
	return len(probesDatas), monitorErr, nil
}

func (c *Controller) SendMetrics(metrics metrics.IMetrics) error {
	return c.pushGateway.Send(metrics.GetId(), metrics.GetMetrics())
}
