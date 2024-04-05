package main

import (
	"context"
	"fmt"
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
	CLOUDWATCHPREFIX string = "CLOUDWATCHPREFIX"
	PUSHGATEWAY      string = "PUSHGATEWAY"
	SQSQUEUENAME     string = "SQS_QUEUE_NAME"
	AWSREGIONCENTRAL string = "AWS_REGION_CENTRAL"
)

func init() {
	logger.SetupLog()

	pushGateway := os.Getenv(PUSHGATEWAY)
	if pushGateway == "" {
		slog.Info(fmt.Sprintf("%s not set, metrics will not be pushed to prometheus", PUSHGATEWAY))
	}

	cloudWatchPrefix := os.Getenv(CLOUDWATCHPREFIX)
	if cloudWatchPrefix == "" {
		slog.Info(fmt.Sprintf("%v not set, metrics will not be pushed to cloudwatch", CLOUDWATCHPREFIX))
	}

	sqsQueueName := os.Getenv(SQSQUEUENAME)
	if sqsQueueName == "" {
		slog.Info(fmt.Sprintf("%v not set, status will not be pushed", SQSQUEUENAME))
	}

	region := os.Getenv(AWSREGIONCENTRAL)
	if region == "" {
		slog.Error(fmt.Sprintf("%v not set, exiting", AWSREGIONCENTRAL))
		os.Exit(1)
	}

	c, err = controller.NewController(ctx, region, pushGateway, sqsQueueName, cloudWatchPrefix)
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
