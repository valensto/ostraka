package provider

import (
	"encoding/json"
	"fmt"
	"github.com/valensto/ostraka/internal/workflow"
	"github.com/valensto/ostraka/internal/workflow/middleware"
	"github.com/valensto/ostraka/internal/workflow/provider/mqtt"
	"github.com/valensto/ostraka/internal/workflow/provider/sse"
	"github.com/valensto/ostraka/internal/workflow/provider/webhook"
)

func NewPublisher(output *workflow.Output, params any, middlewares *middleware.Middlewares) (workflow.Publisher, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling input params: %w", err)
	}

	switch output.Destination {
	case sse.SSE:
		return sse.NewPublisher(output, b, middlewares)

	case mqtt.MQTT:
		return mqtt.NewPublisher(output, b)

	default:
		return nil, fmt.Errorf("unknown publisher type: %s", output.Destination)
	}
}

func NewSubscriber(input *workflow.Input, params any, middlewares *middleware.Middlewares) (workflow.Subscriber, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling input params: %w", err)
	}

	switch input.Source {
	case webhook.Webhook:
		return webhook.NewSubscriber(input, b, middlewares)

	case mqtt.MQTT:
		return mqtt.NewSubscriber(input, b)

	default:
		return nil, fmt.Errorf("unknown subscriber type: %s", input.Source)
	}
}
