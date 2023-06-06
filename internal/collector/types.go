package collector

import (
	"encoding/json"
	"github.com/valensto/ostraka/internal/logger"
	"time"
)

type (
	state string
)

const (
	succeed state = "succeed"
	failed  state = "failed"
)

type event struct {
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

func (e *event) marshall() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		logger.Get().Error().Msgf("error %s marshalling event: %+v", err.Error(), e)
		return nil
	}

	return marshal
}
