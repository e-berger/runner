package controller

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/e-berger/sheepdog-runner/internal/infra"
	"github.com/e-berger/sheepdog-runner/internal/infra/messaging"
	"github.com/e-berger/sheepdog-runner/internal/metrics"
	"github.com/e-berger/sheepdog-runner/internal/probes"
	"github.com/e-berger/sheepdog-runner/internal/status"
)

type Controller struct {
	pushGateway    *metrics.Push
	queueMessaging *messaging.Messaging
	ctx            context.Context
}

func NewController(ctx context.Context, pushGateway string, sqsQueueName string) (*Controller, error) {
	var p *metrics.Push
	if pushGateway != "" {
		p = metrics.NewPush(pushGateway)
	}

	var m *messaging.Messaging
	if sqsQueueName != "" {
		slog.Info("Using messaging", "queue", sqsQueueName)
		cfg, err := infra.NewSession()
		if err != nil {
			return nil, err
		}
		clientSqs := sqs.NewFromConfig(*cfg)
		m = messaging.NewMessaging(clientSqs, sqsQueueName)
		if err := m.Start(ctx); err != nil {
			return nil, err
		}
	}

	return &Controller{
		pushGateway:    p,
		queueMessaging: m,
		ctx:            ctx,
	}, nil
}

func (c *Controller) Run(probesDatas []probes.IProbe) {
	monitorErr := 0
	wg := new(sync.WaitGroup)
	for _, probe := range probesDatas {
		wg.Add(1)
		go c.runProbe(probe, wg, &monitorErr)
	}
	wg.Wait()
	slog.Info("End monitoring", "nb error", monitorErr)
}

func (c *Controller) runProbe(probe probes.IProbe, wg *sync.WaitGroup, monitorErr *int) {
	slog.Info("Launching monitoring", "probe", probe.String())

	defer wg.Done()
	// Launch the probe
	result, errProbe := probe.Launch()
	if errProbe != nil {
		*monitorErr++
	} else {
		// Push metrics to PushGateway
		err := c.SendMetrics(result)
		if err != nil {
			*monitorErr++
			slog.Error("Error pushing monitoring", "error", err)
		}
	}
	// Find out if we need to send an update for status
	if c.queueMessaging != nil {
		errStatus := c.UpdateProbeStatus(probe, result.GetTime(), errProbe)
		if errStatus != nil {
			slog.Error("Error publishing status", "error", errStatus)
		}
	} else {
		slog.Info("No queue messaging defined")
	}
}

func (c *Controller) SendMetrics(metrics metrics.IMetrics) error {
	var err error
	if c.pushGateway != nil {
		slog.Info("Metrics monitoring", "probe", metrics.String())
		err = c.pushGateway.Send(metrics.GetId(), metrics.GetMetrics())
		if err != nil {
			slog.Error("Error pushing monitoring", "error", err)
		}
	} else {
		slog.Info("No push gateway defined")
	}
	return err
}

func (c *Controller) UpdateProbeStatus(probe probes.IProbe, started time.Time, err error) error {
	// Detect if probe status has changed
	if err != nil || (err == nil && probe.IsError()) {
		slog.Info("Update status", "probe", probe.GetId(), "current error", probe.IsError(), "new error", err)
		var s *status.Status
		if err != nil {
			s = status.NewStatus(started, probe.GetId(), probes.ERROR, err.Error(), probe.GetMode(), probe.GetLocation())
		} else {
			s = status.NewStatus(started, probe.GetId(), probes.UP, "", probe.GetMode(), probe.GetLocation())
		}
		return c.queueMessaging.Publish(c.ctx, s)
	} else {
		slog.Debug("No status update", "probe", probe.GetId(), "current error", probe.IsError())
	}
	return nil
}
