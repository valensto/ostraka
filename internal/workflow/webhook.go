package workflow

import "fmt"

type WebhookParams struct {
	Endpoint string `json:"endpoint"`
}

func (w WebhookParams) validate() error {
	if w.Endpoint == "" {
		return fmt.Errorf("webhook endpoint is empty")
	}

	return nil
}

func (i *Input) WebhookParams() (WebhookParams, error) {
	if i.Source != Webhook {
		return WebhookParams{}, fmt.Errorf("input source is not Webhook")
	}

	params, ok := i.params.(WebhookParams)
	if !ok {
		return WebhookParams{}, fmt.Errorf("input params are not of type WebhookParams")
	}

	return params, nil
}
