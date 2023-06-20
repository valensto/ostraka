package workflow

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/valensto/ostraka/internal/extractor"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/provider"
)

type Extractor interface {
	Extract() (extractor.Files, error)
}

type Builder struct {
	server    *http.Server
	consumers []consumer
	extractor Extractor
}

func NewBuilder(extractor Extractor, server *http.Server, consumers ...consumer) *Builder {
	return &Builder{
		server:    server,
		consumers: consumers,
		extractor: extractor,
	}
}

func (b *Builder) Build() ([]*Workflow, error) {
	validate := validator.New()

	files, err := b.extractor.Extract()
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

		for s, cors := range wf.Middlewares.CORS {
			fmt.Println(s, cors)
		}

		wf.consumers = b.consumers
		opts := provider.Options{
			Middlewares: wf.Middlewares,
			Server:      b.server,
		}

		if err := wf.Init(opts); err != nil {
			return nil, fmt.Errorf("error initializing workflow %s: %w", ext, err)
		}
	}

	return workflows, nil
}
