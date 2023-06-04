package collector

import (
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
)

type consumer interface {
	Consume(bytes []byte)
}

type notifier interface {
	FullName() string
}

type Collector struct {
	WorkflowSlug string
	Consumers    []consumer
	Queue        chan Event
}

func New(workflowSlug string, consumers ...consumer) *Collector {
	c := &Collector{
		WorkflowSlug: workflowSlug,
		Queue:        make(chan Event),
		Consumers:    consumers,
	}

	c.broadcast()
	return c
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

func (c *Collector) Collect(notifier notifier, data []byte, err error) {
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
		WorkflowSlug: c.WorkflowSlug,
		Action:       nAction,
		Notifier:     notifier.FullName(),
		Data:         string(data),
		State:        nStatus,
		Message:      "message",
	}

	c.Queue <- n
}
