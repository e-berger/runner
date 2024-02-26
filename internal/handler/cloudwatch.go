package handler

import (
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/e-berger/sheepdog-runner/internal/controller"
)

func CloudWatchEventHandler(c *controller.Controller, cloudWatchEvent events.CloudWatchEvent) (Response, error) {
	var response Response
	slog.Info("CloudWatch Event")
	return response, nil
}
