package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/teramono/utilities/pkg/setup"
)

// APIServer ...
type APIServer struct {
	setup setup.Setup
}

// NewAPIServer ...
func NewAPIServer(setup setup.Setup) (APIServer, error) {
	return APIServer{
		setup: setup,
	}, nil
}

// Listen ...
func (server *APIServer) Listen() error {
	router := gin.Default()

	// Use GRPC instead.
	router.POST("/", func(c *gin.Context) {})

	return router.Run()
}
