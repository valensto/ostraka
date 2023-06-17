package workflow

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/valensto/ostraka/internal/logger"
	"time"
)

type consumer interface {
	Consume(event Event)
}

type Collector struct {
	workflow  *Workflow
	consumers []consumer
	queue     chan Event
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
	event *EventCollect
	err   error

	queue    chan<- Event
	logLevel zerolog.Level
}

func (c *Collector) Collect(from *Input, data []byte) *Collect {
	return &Collect{
		event: &EventCollect{
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

func (c *Collect) WithOutput(output *Output, event Event) *Collect {
	c.event.To = source{
		Provider: output.Destination,
		Name:     output.Name,
		Data:     "string(event.jsonEncode())",
	}
	return c
}

func (c *Collect) WithError(err error) *Collect {
	if err == nil {
		return c
	}

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
	/*c.queue <- *c.event*/
}

func (c *Collect) Error() error {
	return c.err
}

type (
	state string
)

const (
	succeed state = "succeed"
	failed  state = "failed"
)

type EventCollect struct {
	Id           string    `json:"id"`
	WorkflowSlug string    `json:"workflow_slug"`
	From         source    `json:"from"`
	To           source    `json:"to"`
	State        state     `json:"state"`
	Message      string    `json:"message"`
	CollectedAt  time.Time `json:"collected_at"`
}

type source struct {
	Provider string `json:"provider"`
	Name     string `json:"name"`
	Data     string `json:"data"`
}

func (e *EventCollect) ToMap() map[string]any {
	return map[string]any{
		"id":            e.Id,
		"workflow_slug": e.WorkflowSlug,
		"from": map[string]any{
			"provider": e.From.Provider,
			"name":     e.From.Name,
			"data":     e.From.Data,
		},
		"to": map[string]any{
			"provider": e.To.Provider,
			"name":     e.To.Name,
			"data":     e.To.Data,
		},
		"state":        e.State,
		"message":      e.Message,
		"collected_at": e.CollectedAt,
	}
}
