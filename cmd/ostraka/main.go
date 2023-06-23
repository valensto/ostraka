package main

import (
	"fmt"
	"github.com/valensto/ostraka/internal/config/static/local"
	"github.com/valensto/ostraka/internal/env"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/webui"
	"github.com/valensto/ostraka/internal/workflow"
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
	localProvider := local.New(".ostraka/workflows")

	ui, err := webui.New(config.Webui, server)
	if err != nil {
		return fmt.Errorf("cannot create webui: %w", err)
	}

	builder := workflow.NewBuilder(localProvider, server, ui)
	workflows, err := builder.Build()
	if err != nil {
		return err
	}

	ui.Serve(workflows)
	return server.Serve()
}
