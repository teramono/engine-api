package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/teramono/utilities/pkg/broker"
	"github.com/teramono/utilities/pkg/setup"
)

// APIServer ...
type APIServer struct {
	setup.APIEngineSetup
	backends []broker.Address
}

// NewAPIServer ...
func NewAPIServer(setup setup.APIEngineSetup) (APIServer, error) {
	return APIServer{
		APIEngineSetup: setup,
		backends:       []broker.Address{},
	}, nil
}

// Listen ...
func (server *APIServer) Listen() error {
	router := gin.Default()

	// API-specific routes
	router.GET("/", server.LoadGigamonoFramework)
	workpacesRouter := router.Group("/workspaces")
	{
		workpacesRouter.POST("", server.CreateWorkspaces)
		workpacesRouter.GET("", server.GetWorkspaces)
	}

	// Backend-specific routes
	router.POST("/login", server.Login)
	router.Any("/run/*all", server.Run)

	return router.Run(fmt.Sprintf(":%d", server.Config.Engines.API.Port))
}
