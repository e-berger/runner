package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	b64 "encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/e-berger/sheepdog-domain/types"
	"github.com/e-berger/sheepdog-runner/internal/controller"
	"github.com/e-berger/sheepdog-runner/internal/probes"
)

func NewAPIGatewayProxyResponse(statusCode int, body string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}

func ApiGatewayEventHandler(c *controller.Controller, apiEvent events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	slog.Info("Apigateway Event", "event", apiEvent)
	var body = apiEvent.Body
	if apiEvent.IsBase64Encoded {
		slog.Debug("Event is base64 encoded")
		bodyDecoded, err := b64.StdEncoding.DecodeString(body)
		if err != nil {
			return NewAPIGatewayProxyResponse(http.StatusBadRequest, "Bad Request, could not decode base64")
		}
		body = string(bodyDecoded)
	}

	var event EventProbes
	if err := json.Unmarshal([]byte(body), &event); err != nil {
		return NewAPIGatewayProxyResponse(http.StatusBadRequest, "Bad Request, could not unmarshal event")
	}

	location, err := types.ParseLocation(event.Location)
	if err != nil {
		return NewAPIGatewayProxyResponse(http.StatusBadRequest, "Bad Request, could not unmarshal location")
	}

	mode, err := types.ParseMode(event.Mode)
	if err != nil {
		return NewAPIGatewayProxyResponse(http.StatusBadRequest, "Bad Request, could not unmarshal mode")
	}

	probeList := probes.Probes{
		Location: location,
		Mode:     mode,
	}

	for _, item := range event.Items {
		probeJSON := &probes.ProbeJSON{}
		err := json.Unmarshal(item, probeJSON)
		if err != nil {
			return NewAPIGatewayProxyResponse(http.StatusBadRequest, "Bad Request, could not unmarshal items")
		}
		probe, err := probes.NewProbeFromJSON(*probeJSON, event.Location)
		if err != nil {
			return NewAPIGatewayProxyResponse(http.StatusBadRequest, "Bad Request, could not unmarshal probes")
		}
		probeList.Probes = append(probeList.Probes, probe)
	}
	results := c.Run(probeList)
	response, err := json.Marshal(results)
	if err != nil {
		return NewAPIGatewayProxyResponse(http.StatusInternalServerError, "Internal Server Error, could not marshal results")
	}
	return NewAPIGatewayProxyResponse(http.StatusOK, string(response))
}
