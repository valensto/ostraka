package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valensto/ostraka/internal/dispatcher"
	"github.com/valensto/ostraka/internal/workflow"
)

func main() {
	port := "4000"
	banner(port)
	if err := run(port); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func run(port string) error {
	workflows, err := workflow.Build()
	if err != nil {
		return err
	}

	return dispatcher.Dispatch(workflows, port)
}

func banner(port string) {
	b := `
 ██████╗ ███████╗████████╗██████╗  █████╗ ██╗  ██╗ █████╗ 
██╔═══██╗██╔════╝╚══██╔══╝██╔══██╗██╔══██╗██║ ██╔╝██╔══██╗
██║   ██║███████╗   ██║   ██████╔╝███████║█████╔╝ ███████║
██║   ██║╚════██║   ██║   ██╔══██╗██╔══██║██╔═██╗ ██╔══██║
╚██████╔╝███████║   ██║   ██║  ██║██║  ██║██║  ██╗██║  ██║
 ╚═════╝ ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝
https://github.com/valensto/ostraka - %v ©
App running on port - %v

`
	t := time.Now()
	y := t.Year()
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	_ = c.Run()
	log.Printf(b, y, port)
}
