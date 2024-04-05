package controller

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/e-berger/sheepdog-domain/status"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/e-berger/sheepdog-runner/internal/infra/messaging"
	"github.com/e-berger/sheepdog-runner/internal/metrics"
	"github.com/e-berger/sheepdog-runner/internal/probes"
	"github.com/e-berger/sheepdog-utils/aws/creds"
	"github.com/e-berger/sheepdog-utils/cfg"
	"github.com/e-berger/sheepdog-utils/cfg/envs"
)

type Controller struct {
	pushGateway    *metrics.Publish
	queueMessaging *messaging.Messaging
	ctx            context.Context
}

func NewController(ctx context.Context, region string, pushGateway string, sqsQueueName string, cloudWatchPrefix string) (*Controller, error) {

	configuration := cfg.NewConfiguration(envs.WithEnvironmentVariables())
	cfg, err := creds.NewSessionForRegion(region, configuration)
	if err != nil {
		return nil, err
	}

	var cw *cloudwatch.Client
	if cloudWatchPrefix != "" {
		cw = cloudwatch.NewFromConfig(*cfg)
	}

	p := metrics.NewPublish(pushGateway, cloudWatchPrefix, cw)

	var m *messaging.Messaging
	if sqsQueueName != "" {
		slog.Info("Using messaging", "queue", sqsQueueName)
		clientSqs := sqs.NewFromConfig(*cfg)
		m = messaging.NewMessaging(clientSqs, sqsQueueName)
		if err := m.Start(ctx); err != nil {
			return nil, err
		}
	}

	return &Controller{
		ctx:            ctx,
		pushGateway:    p,
		queueMessaging: m,
	}, nil
}

func (c *Controller) Run(probesList probes.Probes) {
	monitorErr := 0
	wg := new(sync.WaitGroup)
	for _, probe := range probesList.Probes {
		wg.Add(1)
		go func() {
			slog.Info("Launching monitoring", "probe", probe.String())
			defer wg.Done()
			// Launch the probe
			result, errProbe := probe.Launch()
			if errProbe != nil {
				monitorErr++
			} else {
				// Push metrics
				err := c.SendMetrics(result)
				if err != nil {
					monitorErr++
					slog.Error("Error pushing monitoring", "error", err)
				}
			}
			// Find out if we need to send an update for status
			if c.queueMessaging != nil {
				errStatus := c.UpdateProbeStatus(probe, probesList.Location, probesList.Mode, result.GetTime(), errProbe)
				if errStatus != nil {
					slog.Error("Error publishing status", "error", errStatus)
				}
			} else {
				slog.Info("No queue messaging defined")
			}
		}()
	}
	wg.Wait()
	slog.Info("End monitoring", "nb error", monitorErr)
}

func (c *Controller) SendMetrics(metrics metrics.IMetrics) error {
	slog.Info("Metrics monitoring", "probe", metrics.String())
	if err := c.pushGateway.Send(metrics); err != nil {
		slog.Error("Error pushing monitoring", "error", err)
		return err
	}
	return nil
}

func (c *Controller) UpdateProbeStatus(probe probes.IProbe, location types.Location, mode types.Mode, started time.Time, err error) error {
	if err != nil || probe.IsInError() {
		slog.Info("Update status", "probe", probe.GetId(), "current error", probe.IsInError(), "new error", err)
		var s status.StatusJSON
		if err != nil {
			s = status.NewStatusJSON(started, probe.GetId(), uint(types.ERROR), err.Error(), uint(mode), uint(location))
		} else {
			s = status.NewStatusJSON(started, probe.GetId(), uint(types.UP), "", uint(mode), uint(location))
		}
		return c.queueMessaging.Publish(c.ctx, s)
	} else {
		slog.Debug("No status update", "probe", probe.GetId(), "current error", probe.IsInError())
	}
	return nil
}
