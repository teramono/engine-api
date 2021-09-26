package engine

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/teramono/utilities/pkg/database/models"
	"github.com/teramono/utilities/pkg/messages"
	"github.com/teramono/utilities/pkg/request"
	"github.com/teramono/utilities/pkg/response"
)

// Index ...
func (server *APIServer) Index(ctx *gin.Context) {
}

// Login ...
func (server *APIServer) Login(ctx *gin.Context) {
	worskspaceName := ctx.GetHeader(request.WorkspaceNameHeader)
	if worskspaceName == "" {
		response.SetErrorResponseWithSourceType(
			ctx,
			response.SourceAPIServer,
			messages.InvalidWorkspaceNameHeader,
		)
		return
	}

	workspace := models.Workspace{Name: worskspaceName}
	workspace, err := workspace.GetByName(&server.setup.WorkspacesDB)
	if err != nil {
		response.SetErrorResponseWithSourceType(
			ctx,
			response.SourceAPIServer,
			messages.ErrorMessage(err.Error()), // TODO: User friendly message
		)
		return
	}

	workspaceID := strconv.Itoa(int(workspace.ID))
	if err = server.dispatchViaInterface(ctx, workspaceID, server.backend.Login); err != nil {
		response.SetErrorResponseWithSourceType(
			ctx,
			response.SourceAPIServer,
			messages.ErrorMessage(err.Error()), // TODO: User friendly message
		)
		return
	}
}

// Run ...
func (server *APIServer) Run(ctx *gin.Context) {
	workspaceID := ctx.GetHeader(request.WorkspaceIDHeader)
	if workspaceID == "" {
		response.SetErrorResponseWithSourceType(
			ctx,
			response.SourceAPIServer,
			messages.InvalidWorkspaceIDHeader,
		)
		return
	}

	if err := server.dispatchViaInterface(ctx, workspaceID, server.backend.Login); err != nil {
		response.SetErrorResponseWithSourceType(
			ctx,
			response.SourceAPIServer,
			messages.ErrorMessage(err.Error()), // TODO: User friendly message
		)
		return
	}
}

// dispatchWithInterface
func (server *APIServer) dispatchViaInterface(
	ctx *gin.Context,
	workspaceID string,
	interfaceFunc func(s string, r request.Request) (*http.Response, error),
) error {
	// Construct request.
	req, err := request.NewRequestFromContext(ctx)
	if err != nil {
		return err
	}

	// Call interface function with request.
	resp, err := interfaceFunc(workspaceID, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Get body bytes from response.
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	body := map[string]interface{}{}
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return err
	}

	// Send response from interface.
	ctx.JSON(resp.StatusCode, gin.H{
		"source": response.SourceBackendServer,
		"body":   body,
	})

	return nil
}
