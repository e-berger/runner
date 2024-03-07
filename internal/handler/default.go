package handler

import (
	"log/slog"

	"github.com/e-berger/sheepdog-runner/internal/controller"
)

func DefaultEventHandler(c *controller.Controller) (Response, error) {
	var response Response
	slog.Info("Default Event")
	return response, nil
}
