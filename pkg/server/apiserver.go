package server

import (
	"fmt"

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

func (server *APIServer) Listen() error {
	router := gin.Default()

	// UI routes.
	router.Use(static.Serve("/", static.LocalFile(server.Config.UI.Dir, true)))

	// Workspaces routes.
	workpacesRouter := router.Group("/workspaces")
	{
		workpacesRouter.POST("", server.CreateWorkspaces)
		workpacesRouter.POST("/login", server.LoginWorkspaces)
		workpacesRouter.Any("/run/*all", server.RunWorkspaces)
	}

	return router.Run(fmt.Sprintf(":%d", server.Config.Engines.API.Port))
}
