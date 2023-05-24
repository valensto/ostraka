package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/dispatcher"
	"log"
	"net/http"
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

	router := chi.NewRouter()

	err = dispatcher.Dispatch(conf, router)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":4000", router)
}
