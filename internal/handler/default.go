package handler

import (
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/e-berger/sheepdog-runner/internal/controller"
)

func DefaultEventHandler(c *controller.Controller, event Event) *events.APIGatewayProxyResponse {
	slog.Info("Default Event", "event", event)
	return nil
}
