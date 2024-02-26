package handler

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/e-berger/sheepdog-runner/internal/controller"
)

type Response struct {
	*events.CloudWatchEvent         `json:",omitempty"`
	*events.APIGatewayProxyResponse `json:",omitempty"`
}

type Event struct {
	events.CloudWatchEvent
	events.APIGatewayProxyRequest
}

func (e Event) Handler(c *controller.Controller) (Response, error) {
	switch {
	case e.APIGatewayProxyRequest.Body != "":
		return ApiGatewayEventHandler(c, e.APIGatewayProxyRequest)
	case len(e.CloudWatchEvent.Detail) > 0:
		return CloudWatchEventHandler(c, e.CloudWatchEvent)
	default:
		return DefaultEventHandler(c)
	}
}
