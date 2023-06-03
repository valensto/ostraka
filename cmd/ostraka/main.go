package main

import (
	"context"
	"github.com/valensto/ostraka/internal/config/static/local"
	"github.com/valensto/ostraka/internal/dispatcher"
	"github.com/valensto/ostraka/internal/logger"
)

func main() {
	port := "4000"
	logger.Banner(port)

	if err := run(port); err != nil {
		logger.Get().Fatal().Msg(err.Error())
	}
}

func run(port string) error {
	ctx := context.Background()
	extractor := local.New(".ostraka/workflows")

	return dispatcher.Dispatch(ctx, extractor, port)
}
