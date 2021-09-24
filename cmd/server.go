package main

import (
	"log"

	"github.com/teramono/engine-api/pkg/engine"
	"github.com/teramono/utilities/pkg/setup"
)

func main() {
	// ...
	setup, err := setup.NewSetup()
	if err != nil {
		log.Fatalln(err)
	}

	// ...
	server, err := engine.NewAPIServer(setup)
	if err != nil {
		log.Fatalln(err)
	}

	server.Listen()
}
