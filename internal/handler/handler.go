package handler

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/e-berger/sheepdog-runner/internal/controller"
)

type Event struct {
	events.CloudWatchEvent
	events.APIGatewayProxyRequest
}

func (e Event) Handler(c *controller.Controller) (*events.APIGatewayProxyResponse, error) {
	switch {
	//api gateway
	case e.APIGatewayProxyRequest.Body != "":
		return ApiGatewayEventHandler(c, e.APIGatewayProxyRequest), nil
	// eventbridge / cloudwatch
	case len(e.CloudWatchEvent.Detail) > 0:
		if err := CloudWatchEventHandler(c, e.CloudWatchEvent); err != nil {
			panic(err)
		}
	default:
		return DefaultEventHandler(c, e), nil
	}
	return nil, nil
}
