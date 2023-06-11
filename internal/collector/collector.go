package collector

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
	"time"
)

type consumer interface {
	Consume(event Event)
}

type Collector struct {
	workflow  *workflow.Workflow
	consumers []consumer
	queue     chan Event
}

func New(wf *workflow.Workflow, consumers ...consumer) *Collector {
	c := &Collector{
		workflow:  wf,
		consumers: consumers,
		queue:     make(chan Event),
	}

	c.broadcast()
	return c
}

func (c *Collector) broadcast() {
	if len(c.consumers) == 0 {
		return
	}

	go func() {
		for {
			select {
			case event := <-c.queue:
				for _, consumer := range c.consumers {
					consumer.Consume(event)
				}
			}
		}
	}()
}

type Collect struct {
	event *Event
	err   error

	queue    chan<- Event
	logLevel zerolog.Level
}

func (c *Collector) Collect(from *workflow.Input, data []byte) *Collect {
	return &Collect{
		event: &Event{
			WorkflowSlug: c.workflow.Slug,
			From: source{
				Provider: from.Source,
				Name:     from.Name,
				Data:     string(data),
			},
			State:       succeed,
			CollectedAt: time.Now().UTC(),
			Message:     "event sent successfully",
		},

		logLevel: zerolog.InfoLevel,
		queue:    c.queue,
	}
}

func (c *Collect) WithOutput(output *workflow.Output, event workflow.Event) *Collect {
	c.event.To = source{
		Provider: output.Destination,
		Name:     output.Name,
		Data:     string(event.Bytes()),
	}
	return c
}

func (c *Collect) WithError(err error) *Collect {
	c.err = err
	c.event.Message = err.Error()
	c.event.State = failed
	c.logLevel = zerolog.ErrorLevel
	return c
}

func (c *Collect) WithLogLevel(lvl zerolog.Level) *Collect {
	c.logLevel = lvl
	return c
}

func (c *Collect) Send() {
	c.event.Id = uuid.NewString()
	logger.Get().WithLevel(c.logLevel).Msgf(c.event.Message)
	c.queue <- *c.event
}

func (c *Collect) Error() error {
	return c.err
}
