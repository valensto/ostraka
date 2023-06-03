package collector

import (
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
)

type Consumer interface {
	Consume(bytes []byte)
}

type Notifier interface {
	FullName() string
}

type Collector struct {
	WorkflowName string
	Consumers    []Consumer
	Queue        chan Event
}

func New(workflowName string, consumers ...Consumer) *Collector {
	c := &Collector{
		WorkflowName: workflowName,
		Queue:        make(chan Event),
		Consumers:    consumers,
	}

	c.broadcast()
	return c
}

func (c *Collector) Collect(notifier Notifier, data []byte, err error) {
	if len(c.Consumers) == 0 {
		return
	}

	nStatus := succeed
	if err != nil {
		nStatus = failed
	}

	var nAction action
	switch notifier.(type) {
	case *workflow.Output:
		nAction = sent
	case *workflow.Input:
		nAction = received
	default:
		logger.Get().Error().Msgf("unknown notifier: %+v cannot notify", notifier)
		return
	}

	n := Event{
		WorkflowName: c.WorkflowName,
		Action:       nAction,
		Notifier:     notifier.FullName(),
		Data:         data,
		State:        nStatus,
		Message:      "message",
	}

	c.Queue <- n
}

func (c *Collector) broadcast() {
	if len(c.Consumers) == 0 {
		return
	}

	go func() {
		for {
			select {
			case event := <-c.Queue:
				for _, consumer := range c.Consumers {
					consumer.Consume(event.marshall())
				}
			}
		}
	}()
}
