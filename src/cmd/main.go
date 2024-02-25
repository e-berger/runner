package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/e-berger/sheepdog-runner/src/internal/controller"
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

func parsedFormData(request events.APIGatewayProxyRequest) (int, int, error) {
	parsedFormData, err := url.ParseQuery(request.Body)
	if err != nil {
		return 0, 0, err
	}

	limitParam := parsedFormData.Get("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	offsetParam := parsedFormData.Get("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		offset = 0
	}

	return limit, offset, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	limit, offset, err := parsedFormData(request)
	if err != nil {
		slog.Error("Error parsing form data", "error", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	total, numError, err := c.Run(limit, offset)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	if numError > 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusAccepted,
			Body:       fmt.Sprintf("%d/%d OK", numError, total),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(handler)
}
