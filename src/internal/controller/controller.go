package controller

import (
	"log/slog"
	"sync"

	db "github.com/e-berger/sheepdog-runner/src/internal/database"
	"github.com/e-berger/sheepdog-runner/src/internal/domain"
	"github.com/e-berger/sheepdog-runner/src/internal/metrics"
)

type Controller struct {
	database    *db.TursoDatabase
	pushGateway *metrics.Push
}

func NewController(database string, authToken string, pushGateway string) *Controller {
	var p *metrics.Push
	if pushGateway != "" {
		p = metrics.NewPush(pushGateway)
	}
	return &Controller{
		database:    db.NewTursoDatabase(database, authToken),
		pushGateway: p,
	}
}

func (c *Controller) Run(limit int, offset int) (int, int, error) {
	probesDatas, err := c.database.GetProbes(limit, offset)
	if err != nil {
		slog.Error("Error fetching datas", "error", err)
		return 0, 0, err
	}

	monitorErr := 0
	wg := new(sync.WaitGroup)
	for _, probe := range probesDatas {
		wg.Add(1)
		go func(probe domain.IProbe) {
			slog.Info("Launching monitoring", "probe", probe.String())
			defer wg.Done()
			result, err := probe.Launch()
			if err != nil {
				monitorErr++
				slog.Error("Error launching monitoring", "error", err)
			}
			slog.Info("Metrics monitoring", "probe", probe.GetResult().String())
			if c.pushGateway != nil {
				err = c.pushGateway.Send(result.GetId(), result.GetMetrics())
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
