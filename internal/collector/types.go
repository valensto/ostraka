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

type Event struct {
	WorkflowName string `json:"workflow_name"`
	Action       action `json:"action"`
	Notifier     string `json:"notifier"`
	Data         []byte `json:"data"`
	State        state  `json:"state"`
	Message      string `json:"message"`
}

func (n Event) marshall() []byte {
	marshal, err := json.Marshal(n)
	if err != nil {
		logger.Get().Error().Msgf("error %s marshalling event: %+v", err.Error(), n)
		return nil
	}

	return marshal
}
