package workflow

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type consumer interface {
	Consume(entry Entry)
}

type collector struct {
	consumers []consumer
	buffer    Entry
	entries   []Entry
}

func (c *collector) consumes() {
	if len(c.consumers) == 0 {
		return
	}

	for _, consumer := range c.consumers {
		for _, entry := range c.entries {
			consumer.Consume(entry)
		}
	}
}

func (wf *Workflow) collect(input *Input, bytes []byte) collector {
	entry := Entry{
		Id:       uuid.NewString(),
		Workflow: wf.Slug,
		From: source{
			Name:     input.Name,
			Provider: input.Subscriber.Provider(),
			Data:     string(bytes),
		},
		State:       succeed,
		CollectedAt: time.Now().UTC(),
		Message:     "event received and decoded successfully",
	}

	return collector{
		consumers: wf.consumers,
		buffer:    entry,
	}
}

func (c *collector) withError(err error) {
	if err != nil {
		c.buffer.State = failed
		c.buffer.Message = err.Error()
	}
}

func (c *collector) addOutput(output *Output, bytes []byte, err error) {
	if err != nil {
		c.withError(err)
	}

	entry := c.buffer
	entry.Id = uuid.NewString()
	entry.To = source{
		Name:     output.Name,
		Provider: output.Publisher.Provider(),
		Data:     string(bytes),
	}

	c.entries = append(c.entries, entry)
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
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Data     string `json:"data"`
}

func (e Entry) JSONEncode() ([]byte, error) {
	b, ok := json.Marshal(e)
	if ok != nil {
		return nil, fmt.Errorf("error marshalling eventType to json: %w", ok)
	}
	return b, nil
}
