package handler

import (
	"log/slog"

	"github.com/e-berger/sheepdog-runner/internal/controller"
)

func DefaultEventHandler(c *controller.Controller, event Event) (Response, error) {
	var response Response
	slog.Info("Default Event", "event", event)
	return response, nil
}
