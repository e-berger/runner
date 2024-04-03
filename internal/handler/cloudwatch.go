package handler

import (
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	domain "github.com/e-berger/sheepdog-domain/probes"
	"github.com/e-berger/sheepdog-runner/internal/controller"
	"github.com/e-berger/sheepdog-runner/internal/probes"
)

type EventProbes struct {
	Location string            `json:"location"`
	Items    []json.RawMessage `json:"items"`
	Mode     string            `json:"mode"`
}

func CloudWatchEventHandler(c *controller.Controller, cloudWatchEvent events.CloudWatchEvent) error {
	slog.Debug("CloudWatch Event", "event", cloudWatchEvent)
	var event EventProbes
	if err := json.Unmarshal(cloudWatchEvent.Detail, &event); err != nil {
		return err
	}

	location, err := domain.ParseLocation(event.Location)
	if err != nil {
		return err
	}

	mode, err := domain.ParseMode(event.Mode)
	if err != nil {
		return err
	}

	probeList := probes.Probes{
		Location: location,
		Mode:     mode,
	}

	for _, item := range event.Items {
		probeJSON := &probes.ProbeJSON{}
		err := json.Unmarshal(item, probeJSON)
		if err != nil {
			return err
		}
		probe, err := probes.NewProbeFromJSON(*probeJSON, event.Location)
		if err != nil {
			return err
		}
		probeList.Probes = append(probeList.Probes, probe)
	}
	c.Run(probeList)
	return nil
}
