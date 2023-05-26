package main

import (
	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/dispatcher"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	port := "4000"
	banner(port)
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	return dispatcher.Dispatch(conf, port)
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
