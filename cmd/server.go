package main

import (
	"log"

	"github.com/teramono/engine-api/pkg/server"
	"github.com/teramono/utilities/pkg/setup"
)

func main() {
	// ...
	setup, err := setup.NewAPIEngineSetup()
	if err != nil {
		log.Fatalln(err)
	}

	// ...
	server, err := server.NewAPIServer(setup)
	if err != nil {
		log.Fatalln(err)
	}

	server.Listen()
}
