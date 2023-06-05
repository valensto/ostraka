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
	notifier Notifier
	data     []byte
	err      error
}

type Collect struct {
	workflowSlug string
	collection   chan collect
	broadcast    chan<- event
}

func (c *Collector) NewCollect() Collect {
	return Collect{
		workflowSlug: c.workflow.Slug,
		collection:   make(chan collect, len(c.workflow.Outputs)+1),
		broadcast:    c.queue,
	}
}

func (c Collect) Add(from Notifier, data []byte, err error, lvl ...zerolog.Level) {
	logger.LogErr(err, lvl...)

	c.collection <- collect{
		notifier: from,
		data:     data,
		err:      err,
	}
}

func (c Collect) Consume() {
	close(c.collection)

	// TODO: refactor this not proud of it
	// this is a hack to get the same uuid for both sent and received events
	// this is needed to link the sent and received events in the webui
	id := uuid.NewString()
	for col := range c.collection {
		a := sent
		if _, ok := col.notifier.(*workflow.Input); ok {
			id = uuid.NewString()
			a = received
		}

		e := event{
			RelatedId:    id,
			WorkflowSlug: c.workflowSlug,
			Action:       a,
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
