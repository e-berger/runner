package handler

import (
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/e-berger/sheepdog-runner/internal/controller"
	"github.com/e-berger/sheepdog-runner/internal/probes"
)

type EventProbes struct {
	Location string            `json:"location"`
	Items    []json.RawMessage `json:"items"`
	Mode     string            `json:"mode,"`
}

func CloudWatchEventHandler(c *controller.Controller, cloudWatchEvent events.CloudWatchEvent) (Response, error) {
	slog.Info("CloudWatch Event", "event", cloudWatchEvent)
	var response = Response{}
	var event EventProbes
	err := json.Unmarshal(cloudWatchEvent.Detail, &event)
	if err != nil {
		return response, err
	}
	var probeDatas []probes.IProbe
	for _, item := range event.Items {
		probe, err := probes.UnmarshalJSON(item, event.Location, event.Mode)
		if err != nil {
			return response, err
		}
		p, err := probe.CreateProbeFromType()
		if err != nil {
			return response, err
		}
		probeDatas = append(probeDatas, p)
	}
	c.Run(probeDatas)
	return response, nil
}
