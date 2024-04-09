package controller

import (
	"context"
	"log/slog"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/e-berger/sheepdog-domain/events"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/e-berger/sheepdog-runner/internal/infra/messaging"
	"github.com/e-berger/sheepdog-runner/internal/probes"
	"github.com/e-berger/sheepdog-runner/internal/results"
	"github.com/e-berger/sheepdog-utils/aws/creds"
	"github.com/e-berger/sheepdog-utils/cfg"
	"github.com/e-berger/sheepdog-utils/cfg/envs"
)

type Controller struct {
	pushGateway    *results.Publish
	queueMessaging *messaging.Messaging
	ctx            context.Context
}

type resultChannel struct {
	result results.IResults
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

	p := results.NewPublish(pushGateway, cloudWatchPrefix, cw)

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

func (c *Controller) Run(probesList probes.Probes) []results.IResults {
	var results []results.IResults
	wg := new(sync.WaitGroup)

	ch := make(chan resultChannel)
	for _, probe := range probesList.Probes {
		wg.Add(1)
		go c.runProbe(ch, wg, probe, probesList.Location, probesList.Mode)
	}

	wg.Add(1)
	go func() {
		for v := range ch {
			results = append(results, v.result)
			if len(results) == len(probesList.Probes) {
				wg.Done()
			}
		}
	}()
	wg.Wait()
	close(ch)
	return results
}

func (c *Controller) runProbe(ch chan resultChannel, wg *sync.WaitGroup, probe probes.IProbe, location types.Location, mode types.Mode) {
	defer wg.Done()
	slog.Info("Launching monitoring", "probe", probe.String())
	// Launch the probe
	result := probe.Launch(probe.GetHttpClient())
	if result.GetErrorProbe() == nil {
		// Push metrics
		err := c.SendResults(result)
		if err != nil {
			slog.Error("Error pushing monitoring", "error", err)
			result.SetError(err)
		}
	}
	// Find out if we need to send an update for status
	if c.queueMessaging != nil && mode != types.MANUAL {
		errStatus := c.UpdateProbeStatus(probe, location, mode, result)
		if errStatus != nil {
			slog.Error("Error publishing status", "error", errStatus)
			result.SetError(errStatus)
		}
	} else {
		slog.Info("No queue for messaging defined")
	}
	ch <- resultChannel{result: result}
}

func (c *Controller) SendResults(result results.IResults) error {
	slog.Info("Metrics monitoring", "probe", result.String())
	if err := c.pushGateway.Send(result); err != nil {
		return err
	}
	return nil
}

func (c *Controller) UpdateProbeStatus(probe probes.IProbe, location types.Location, mode types.Mode, result results.IResults) error {
	// Probe in error or was in error
	if result.GetErrorProbe() != nil || probe.IsInError() {
		slog.Info("Update events", "probe", probe.GetId(), "current error", probe.IsInError(), "new error", result.GetErrorProbe())
		var e events.EventsJSON
		if result.GetErrorProbe() != nil {
			e = events.NewEventsJSON(result.GetTime(), probe.GetId(), uint(result.GetCode()), result.GetErrorProbe().Error(), uint(mode), uint(location))
		} else {
			e = events.NewEventsJSON(result.GetTime(), probe.GetId(), uint(result.GetCode()), "", uint(mode), uint(location))
		}
		return c.queueMessaging.Publish(c.ctx, e)
	} else {
		slog.Debug("No events update", "probe", probe.GetId(), "current error", probe.IsInError())
	}
	return nil
}
