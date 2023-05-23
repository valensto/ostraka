package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/dispatcher"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	d := dispatcher.New(conf, chi.NewRouter())
	return d.Dispatch()
}
