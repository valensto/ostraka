package mqtt

import (
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
)

func NewPublisher(output workflow.Output) (*MQTT, error) {
	params, err := output.MQTTParams()
	if err != nil {
		return nil, err
	}

	c, err := connect(output.Name, params)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *MQTT) Register(events <-chan []byte) error {
	l := logger.Get()
	l.Info().Msgf("output %s of type MQTT registered. Publishing to topic %s", c.name, c.params.Topic)

	go func() {
		for {
			select {
			case event := <-events:
				token := c.client.Publish(c.params.Topic, 1, false, event)
				token.Wait()
				if token.Error() != nil {
					l.Error().Msgf("error publishing to topic: %s", c.params.Topic)
				}
			}
		}
	}()

	return nil
}
