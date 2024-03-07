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
var err error
var ctx = context.Background()

func init() {
	lvl := new(slog.LevelVar)
	logLevel := os.Getenv("LOGLEVEL")
	lvl.Set(slog.LevelInfo)
	if logLevel != "" {
		slog.Info("Logger", "loglevel", logLevel)
		lvl.UnmarshalText([]byte(logLevel))
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)

	pushgateway := os.Getenv("PUSHGATEWAY")
	if pushgateway == "" {
		slog.Info("PUSHGATEWAY not set, metrics will not be pushed")
	}

	sqsQueueName := os.Getenv("SQS_QUEUE_NAME")
	c, err = controller.NewController(ctx, pushgateway, sqsQueueName)
	if err != nil {
		slog.Error("Creating controller", "error", err)
		os.Exit(1)
	}
}

func mainHandler(_ context.Context, event handler.Event) (handler.Response, error) {
	return event.Handler(c)
}

func main() {
	lambda.Start(mainHandler)
}
