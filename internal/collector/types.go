package collector

import (
	"encoding/json"
	"github.com/valensto/ostraka/internal/logger"
)

type (
	state  string
	action string
)

func (s state) String() string {
	return string(s)
}

func (a action) String() string {
	return string(a)
}

const (
	succeed state = "succeed"
	failed  state = "failed"

	received action = "received"
	sent     action = "sent"
)

type event struct {
	RelatedId    string `json:"related_id"`
	WorkflowSlug string `json:"workflow_slug"`
	Action       action `json:"action"`
	Notifier     string `json:"notifier"`
	Provider     string `json:"provider"`
	Data         string `json:"data"`
	State        state  `json:"state"`
	Message      string `json:"message"`
}

func (n event) marshall() []byte {
	marshal, err := json.Marshal(n)
	if err != nil {
		logger.Get().Error().Msgf("error %s marshalling event: %+v", err.Error(), n)
		return nil
	}

	return marshal
}
