package workflow

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/provider"
)

type Builder struct {
	server    *http.Server
	consumers []consumer
	provider  config.Provider
}

func NewBuilder(provider config.Provider, server *http.Server, consumers ...consumer) *Builder {
	return &Builder{
		server:    server,
		consumers: consumers,
		provider:  provider,
	}
}

func (b *Builder) Build() ([]*Workflow, error) {
	validate := validator.New()

	files, err := b.provider.Extract()
	if err != nil {
		return nil, fmt.Errorf("error extracting workflows: %w", err)
	}

	var workflows []*Workflow
	for ext, bytes := range files {
		var wf Workflow
		if err := ext.Unmarshal(bytes, &wf); err != nil {
			return nil, fmt.Errorf("error unmarshalling workflow %s: %w", ext, err)
		}

		if err := validate.Struct(wf); err != nil {
			return nil, fmt.Errorf("error validating workflow %s: %w", ext, err)
		}

		wf.consumers = b.consumers
		opts, err := provider.NewOptions(b.server, wf.Middlewares)
		if err != nil {
			return nil, fmt.Errorf("error initializing workflow %s: %w", ext, err)
		}

		if err := wf.Init(opts); err != nil {
			return nil, fmt.Errorf("error initializing workflow %s: %w", ext, err)
		}

		workflows = append(workflows, &wf)
	}

	return workflows, nil
}
