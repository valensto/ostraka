package collector

import (
	"time"
)

type (
	state string
)

const (
	succeed state = "succeed"
	failed  state = "failed"
)

type Event struct {
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

func (e *Event) ToMap() map[string]any {
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
