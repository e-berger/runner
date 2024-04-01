package handler

import (
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/e-berger/sheepdog-runner/internal/controller"
)

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

func ApiGatewayEventHandler(c *controller.Controller, apiEvent events.APIGatewayProxyRequest) (Response, error) {
	slog.Info("Apigateway Event")
	var response = Response{}
	response.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}
	return response, nil
}
