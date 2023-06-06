package main

import (
	"fmt"

	"github.com/valensto/ostraka/internal/config/env"
	"github.com/valensto/ostraka/internal/config/static/local"
	"github.com/valensto/ostraka/internal/dispatcher"
	"github.com/valensto/ostraka/internal/logger"
)

func main() {
	logger.Banner()
	if err := run(); err != nil {
		logger.Get().Fatal().Msg(err.Error())
	}
}

func run() error {
	config, err := env.Load()
	if err != nil {
		return fmt.Errorf("cannot load env config: %w", err)
	}

	workflows, err := local.Extract(".ostraka/workflows")
	if err != nil {
		return fmt.Errorf("cannot extract workflows: %w", err)
	}

	return dispatcher.Dispatch(config, workflows)
}
