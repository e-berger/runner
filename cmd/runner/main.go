package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/e-berger/sheepdog-runner/internal/controller"
	"github.com/e-berger/sheepdog-runner/internal/handler"
	"github.com/e-berger/sheepdog-runner/internal/infra/logger"
)

var c *controller.Controller
var err error
var ctx = context.Background()

const (
	PUSHGATEWAY  string = "PUSHGATEWAY"
	SQSQUEUENAME string = "SQS_QUEUE_NAME"
)

func init() {
	logger.SetupLog()

	pushgateway := os.Getenv(PUSHGATEWAY)
	if pushgateway == "" {
		slog.Info("PUSHGATEWAY not set, metrics will not be pushed")
	}

	sqsQueueName := os.Getenv(SQSQUEUENAME)
	if sqsQueueName == "" {
		slog.Info("SQS_QUEUE_NAME not set, status will not be pushed")
	}

	c, err = controller.NewController(ctx, pushgateway, sqsQueueName)
	if err != nil {
		slog.Error("Creating controller", "error", err)
		os.Exit(1)
	}
}

func mainHandler(_ context.Context, event handler.Event) (*events.APIGatewayProxyResponse, error) {
	return event.Handler(c)
}

func main() {
	lambda.Start(mainHandler)
}
