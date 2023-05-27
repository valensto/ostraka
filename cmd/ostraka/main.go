package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/dispatcher"
)

func main() {
	if err := run("4000"); err != nil {
		log.Fatal(err)
	}
}

func run(port string) error {
	banner(port)
	workflow, err := config.NewWorkflow()
	if err != nil {
		return err
	}

	return dispatcher.Dispatch(workflow, port)
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
