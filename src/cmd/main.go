package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/e-berger/sheepdog-runner/src/internal/datas"
	"github.com/e-berger/sheepdog-runner/src/internal/monitor"
)

var primaryUrl string
var authToken string

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	database := os.Getenv("TURSO_DATABASE")
	authToken = os.Getenv("TURSO_TOKEN")
	primaryUrl = fmt.Sprintf("https://%s.turso.io", database)
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

	resutls, err := datas.Fetch(limit, offset, primaryUrl, authToken)
	if err != nil {
		slog.Error("Error fetching datas", "error", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	var m monitor.IMonitor
	monitorErr := 0

	wg := new(sync.WaitGroup)
	for _, result := range resutls {
		m, err = monitor.GetMonitoring(result)
		if err != nil {
			monitorErr++
			slog.Error("Error getting monitoring", "error", err)
			continue
		}
		wg.Add(1)
		go func(m monitor.IMonitor) {
			defer wg.Done()
			err = m.Launch()

			if err != nil {
				monitorErr++
				slog.Error("Error launching monitoring", "error", err)
			}
		}(m)
	}
	wg.Wait()

	if monitorErr > 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusAccepted,
			Body:       fmt.Sprintf("%d/%d OK", monitorErr, len(resutls)),
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
