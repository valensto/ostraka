package workflow

import (
	"context"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/valensto/ostraka/internal/server"
)

type Publisher interface {
	Publish(events <-chan Event, mux *server.Server) error
	Output() *Output
}

type Subscriber interface {
	Subscribe(dispatch func(ctx context.Context, input *Input, data []byte) error, mux *server.Server) error
	Input() *Input
}

type Workflow struct {
	Name        string
	Slug        string
	Subscribers []Subscriber
	Publishers  []Publisher
}

func New(name string, subscribers []Subscriber, publishers []Publisher) (*Workflow, error) {
	if name == "" {
		return nil, fmt.Errorf("workflow name is empty")
	}

	return &Workflow{
		Name:        name,
		Slug:        slug.Make(name),
		Subscribers: subscribers,
		Publishers:  publishers,
	}, nil
}
