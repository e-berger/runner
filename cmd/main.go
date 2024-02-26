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
	database := os.Getenv("TURSO_DATABASE")
	if database == "" {
		slog.Error("TURSO_DATABASE not set")
		panic("TURSO_DATABASE not set")
	}
	authToken := os.Getenv("TURSO_TOKEN")
	if authToken == "" {
		slog.Error("TURSO_TOKEN not set")
		panic("TURSO_TOKEN not set")
	}
	pushgateway := os.Getenv("PUSHGATEWAY")
	if pushgateway == "" {
		slog.Info("PUSHGATEWAY not set, metrics will not be pushed")
	}
	c = controller.NewController(database, authToken, pushgateway)
}

func mainHandler(_ context.Context, event handler.Event) (handler.Response, error) {
	return event.Handler(c)
}

func main() {
	lambda.Start(mainHandler)
}
