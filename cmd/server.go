package main

import (
	"log"

	"github.com/teramono/engine-api/pkg/server"
	"github.com/teramono/utilities/pkg/setup"
)

func main() {
	// Establish api engine setup.
	setup, err := setup.NewAPIEngineSetup()
	if err != nil {
		log.Panic(err)
	}

	// Create server.
	server, err := server.NewAPIServer(setup)
	if err != nil {
		log.Panic(err)
	}

	defer server.BrokerClient.Close()

	// Start server.
	server.Listen()
}
