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
	// limit, offset, err := parsedFormData(apiEvent)
	// if err != nil {
	// 	slog.Error("Error parsing form data", "error", err)
	// 	response.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}
	// 	return response, err
	// }

	// probes, err := c.Database.GetProbes(limit, offset)
	// if err != nil {
	// 	slog.Error("Error fetching datas", "error", err)
	// 	response.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}
	// 	return response, err
	// }

	// total, numError, err := c.Run(probes)
	// if err != nil {
	// 	response.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}
	// 	return response, err
	// }

	// if numError > 0 {
	// 	response.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{
	// 		StatusCode: http.StatusAccepted,
	// 		Body:       fmt.Sprintf("%d/%d OK", numError, total),
	// 	}
	// 	return response, nil
	// }
	response.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}
	return response, nil
}
