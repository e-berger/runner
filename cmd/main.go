package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/e-berger/sheepdog-runner/internal/controller"
	"github.com/e-berger/sheepdog-runner/internal/handler"
)

var c *controller.Controller

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	pushgateway := os.Getenv("PUSHGATEWAY")
	if pushgateway == "" {
		slog.Info("PUSHGATEWAY not set, metrics will not be pushed")
	}
	c = controller.NewController(pushgateway)
}

func mainHandler(_ context.Context, event handler.Event) (handler.Response, error) {
	return event.Handler(c)
}

func main() {
	lambda.Start(mainHandler)
}
