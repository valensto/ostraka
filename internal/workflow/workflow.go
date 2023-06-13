package workflow

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow/middleware"
)

type Workflow struct {
	Name    string
	Slug    string
	Inputs  []*Input
	Outputs []*Output
}

type Options struct {
	Middlewares *middleware.Middlewares
	Server      *server.Server
}

func New(name string, inputs []*Input, outputs []*Output) (*Workflow, error) {
	if name == "" {
		return nil, fmt.Errorf("workflow name is empty")
	}

	return &Workflow{
		Name:    name,
		Slug:    slug.Make(name),
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}
