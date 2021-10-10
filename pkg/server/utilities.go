package server

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/teramono/utilities/pkg/database/models"
	"github.com/teramono/utilities/pkg/messages"
	"github.com/teramono/utilities/pkg/request"
)

func (server *APIServer) getHeaders(ctx *gin.Context) map[string][]string {
	headers := map[string][]string{}

	for k, v := range ctx.Request.Header {
		headers[k] = v
	}

	return headers
}

func (server *APIServer) getWorkspaceID(ctx *gin.Context) (uint, error) {
	var workspaceID uint

	// First get workspace id.
	workspaceIDStr := ctx.GetHeader(request.WorkspaceIDHeader)
	if workspaceIDStr == "" {
		// Fallback to using workpace name header if it exists.
		workspaceName := ctx.GetHeader(request.WorkspaceIDHeader)
		if workspaceName == "" {
			return 0, fmt.Errorf(messages.InvalidWorkspaceIDAndNameHeader.String())
		}

		// Get workspace by name.
		workspace, err := (&models.Workspace{Name: workspaceName}).GetByName(&server.DB)
		if err != nil {
			return 0, fmt.Errorf(
				"%s: %w",
				messages.UnableToFindWorkspace(workspaceName).String(),
				err,
			)
		}

		workspaceID = workspace.ID
	} else {
		id, err := strconv.ParseUint(workspaceIDStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf(
				"%s: %w",
				messages.InvalidWorkspaceIDHeader.String(),
				err,
			)
		}

		workspaceID = uint(id)
	}

	return workspaceID, nil
}
