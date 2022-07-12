package main

import (
	"github.com/obgnail/go-spawn/config"
	"github.com/obgnail/go-spawn/ssh"
	"log"
)

func main() {
	cfg := config.Config
	commands := cfg.Commands
	timeout := config.Config.Timeout

	client, err := ssh.Dial(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	spawn := ssh.NewSpawn(session)
	if err := spawn.Interact(commands, timeout); err != nil {
		log.Fatal(err)
	}
}
