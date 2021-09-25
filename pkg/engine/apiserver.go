package engine

import (
	"github.com/gin-gonic/gin"
	backendEngine "github.com/teramono/engine-backend/pkg/engine"
	"github.com/teramono/utilities/pkg/setup"
)

// APIServer ...
type APIServer struct {
	setup setup.Setup
	backend backendEngine.Backend
}

// NewAPIServer ...
func NewAPIServer(setup setup.Setup) (APIServer, error) {
	backend := backendEngine.NewBackend(&setup)
	return APIServer{
		setup: setup,
		backend: backend,
	}, nil
}

// Listen ...
func (server *APIServer) Listen() error {
	router := gin.Default()

	router.POST("/", server.Index) // Serves Gigamono page.
	router.POST("/login", server.Login) // X-WORKSPACE-NAME
	router.POST("/run/*all", server.Run) // X-WORKSPACE-ID

	return router.Run(":5050") // TODO: Get from setup.Config
}
