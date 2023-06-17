package workflow

import (
	"github.com/google/uuid"
	"time"
)

type consumer interface {
	Consume(entry Entry)
}

type collector struct {
	consumers []consumer
	entries   map[*Input]Entry
}

func (c *collector) dump() {
	if len(c.consumers) == 0 {
		return
	}

	for _, consumer := range c.consumers {
		for _, entry := range c.entries {
			go consumer.Consume(entry)
		}
	}
}

func (wf *Workflow) newCollector(from *Input, data []byte) collector {
	col := make(map[*Input]Entry)

	e := Entry{
		Id:       uuid.NewString(),
		Workflow: wf.Slug,
		From: source{
			Provider: "from.Source",
			Name:     from.Name,
			Data:     string(data),
		},
		State:       succeed,
		CollectedAt: time.Now().UTC(),
		Message:     "event sent successfully",
	}

	col[from] = e
	return collector{
		consumers: wf.consumers,
		entries:   col,
	}
}

func (c Entry) send() {

}

type (
	state string
)

const (
	succeed state = "succeed"
	failed  state = "failed"
)

type Entry struct {
	Id          string    `json:"id"`
	Workflow    string    `json:"workflow"`
	From        source    `json:"from"`
	To          source    `json:"to"`
	State       state     `json:"state"`
	Message     string    `json:"message"`
	CollectedAt time.Time `json:"collected_at"`
}

type source struct {
	Provider string `json:"provider"`
	Name     string `json:"name"`
	Data     string `json:"data"`
}
