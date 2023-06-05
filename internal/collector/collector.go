package collector

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
)

type consumer interface {
	Consume(bytes []byte)
}

type Collector struct {
	workflow  *workflow.Workflow
	consumers []consumer
	queue     chan event
}

func New(wf *workflow.Workflow, consumers ...consumer) *Collector {
	c := &Collector{
		workflow:  wf,
		consumers: consumers,
		queue:     make(chan event),
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
					consumer.Consume(event.marshall())
				}
			}
		}
	}()
}

type Notifier interface {
	GetName() string
	GetProvider() string
}

type collect struct {
	relatedId string
	notifier  Notifier
	action    action
	data      []byte
	err       error
}

type Collect struct {
	relatedId    string
	workflowSlug string
	collection   []collect
	broadcast    chan<- event
}

func (c *Collector) NewCollect() Collect {
	return Collect{
		relatedId:    uuid.NewString(),
		workflowSlug: c.workflow.Slug,
		collection:   make([]collect, len(c.workflow.Outputs)+1),
		broadcast:    c.queue,
	}
}

func (c Collect) Add(from Notifier, data []byte, err error, lvl ...zerolog.Level) {
	logger.LogErr(err, lvl...)

	a := sent
	if _, ok := from.(*workflow.Input); ok {
		a = received
	}

	c.collection = append(c.collection, collect{
		relatedId: c.relatedId,
		notifier:  from,
		action:    a,
		data:      data,
		err:       err,
	})

	if a == sent {
		c.relatedId = uuid.NewString()
	}
}

func (c Collect) Consume() {
	for _, col := range c.collection {
		e := event{
			RelatedId:    col.relatedId,
			WorkflowSlug: c.workflowSlug,
			Action:       col.action,
			Notifier:     col.notifier.GetName(),
			Provider:     col.notifier.GetProvider(),
			Data:         string(col.data),
			State:        succeed,
			Message:      "event sent to output",
		}

		if col.err != nil {
			e.State = failed
			e.Message = col.err.Error()
		}

		c.broadcast <- e
	}
}
