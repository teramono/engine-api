package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index ...
func (server *APIServer) LoadGigamonoFramework(ctx *gin.Context) {
}

// Login ...
func (server *APIServer) Login(ctx *gin.Context) {
	// TODO:
	ctx.JSON(http.StatusOK, gin.H{ //
		"message": "Logging user in...",
	})
}

// Run ...
func (server *APIServer) Run(ctx *gin.Context) {
	// TODO:
	ctx.JSON(http.StatusOK, gin.H{ //
		"message": "Script running...",
	})
}

// CreateWorkspaces ...
func (server *APIServer) CreateWorkspaces(ctx *gin.Context) {
	// TODO:
	ctx.JSON(http.StatusOK, gin.H{ //
		"message": "Creating workspaces...",
	})
}

// GetWorkspaces ...
func (server *APIServer) GetWorkspaces(ctx *gin.Context) {
	// TODO:
	ctx.JSON(http.StatusOK, gin.H{ //
		"message": "Getting workspaces...",
	})
}
