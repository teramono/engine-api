package crud

import (
	"github.com/teramono/engine-api/pkg/server"
	"github.com/teramono/utilities/pkg/database/models"
)

func CreateWorkspace(server *server.APIServer, name string) error {
	// Add workspace to db.
	workspace := models.Workspace {
		Name: name,
	}

	return workspace.Create(&server.DB)
}
