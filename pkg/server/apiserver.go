package server

import (
	"fmt"

	"github.com/gin-contrib/location"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/teramono/utilities/pkg/setup"
)

type APIServer struct {
	setup.APIEngineSetup
	Validator *validator.Validate
}

func NewAPIServer(setup setup.APIEngineSetup) (APIServer, error) {
	valid := validator.New()

	return APIServer{
		APIEngineSetup: setup,
		Validator:      valid,
	}, nil
}

func (server *APIServer) LogsVersion() uint {
	return server.Config.Broker.Subscriptions.Logs.Version
}

func (server *APIServer) Listen() error {
	router := gin.Default()

	// Middlewares.
	router.Use(location.Default())

	// UI routes.
	router.Use(static.Serve("/", static.LocalFile(server.Config.UI.Dir, true)))

	// Run route.
	router.Any("/run/*all", server.Run)

	// Workspaces routes.
	workpacesRouter := router.Group("/workspaces")
	{
		workpacesRouter.POST("", server.CreateWorkspaces)
	}

	return router.Run(fmt.Sprintf(":%d", server.Config.Engines.API.Port))
}
