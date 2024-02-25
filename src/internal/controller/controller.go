package controller

import (
	"log/slog"
	"sync"

	db "github.com/e-berger/sheepdog-runner/src/internal/database"
	"github.com/e-berger/sheepdog-runner/src/internal/domain"
	"github.com/e-berger/sheepdog-runner/src/internal/metrics"
)

type Controller struct {
	authToken   string
	database    *db.TursoDatabase
	pushGateway *metrics.Push
}

func NewController(database string, pushGateway string, authToken string) *Controller {
	var p *metrics.Push
	if pushGateway != "" {
		p = metrics.NewPush(pushGateway)
	}
	return &Controller{
		authToken:   authToken,
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
	for _, p := range probesDatas {
		wg.Add(1)
		go func(p domain.IProbe) {
			slog.Info("Launching monitoring", "probe", p.String())
			defer wg.Done()
			err = p.Launch()
			if err != nil {
				monitorErr++
				slog.Error("Error launching monitoring", "error", err)
			}
			slog.Info("Metrics monitoring", "probe", p.GetResult().String())
			if c.pushGateway != nil {
				err = p.Push(c.pushGateway)
				if err != nil {
					monitorErr++
					slog.Error("Error pushing monitoring", "error", err)
				}
			}

		}(p)
	}
	wg.Wait()
	return len(probesDatas), monitorErr, nil
}
