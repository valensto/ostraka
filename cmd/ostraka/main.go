package main

import (
	"github.com/rs/zerolog/log"
	"github.com/valensto/ostraka/internal/dispatcher"
	"github.com/valensto/ostraka/internal/workflow"
	"github.com/valensto/ostraka/logger"
)

func main() {
	port := "4000"
	logger.Banner(port)
	if err := run(port); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func run(port string) error {
	workflows, err := workflow.Build(".ostraka/workflows")
	if err != nil {
		return err
	}

	return dispatcher.Dispatch(workflows, port)
}
