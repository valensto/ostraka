package main

import (
	"fmt"
	"github.com/valensto/ostraka/internal/config/env"
	"github.com/valensto/ostraka/internal/config/static"
	"github.com/valensto/ostraka/internal/config/static/local"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/webui"
)

func main() {
	logger.Banner()
	if err := run(); err != nil {
		logger.Get().Fatal().Msg(err.Error())
	}
}

func run() error {
	config := env.Load()
	server := http.New(config)

	contentFile, err := local.Extract(".ostraka/workflows")
	if err != nil {
		return fmt.Errorf("cannot extract content file: %w", err)
	}

	workflows, err := static.BuildWorkflows(contentFile, server)
	if err != nil {
		return fmt.Errorf("cannot build workflows: %w", err)
	}

	uiConsumer, err := webui.New(config.Webui, server, workflows)
	if err != nil {
		return fmt.Errorf("cannot create webui: %w", err)
	}

	for _, wf := range workflows {
		err := wf.Listen(uiConsumer)
		if err != nil {
			return fmt.Errorf("cannot listen workflow %s: %w", wf.Name, err)
		}
	}

	return server.Run()
}
