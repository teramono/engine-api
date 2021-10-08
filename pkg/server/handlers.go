package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/teramono/utilities/pkg/broker"
	"github.com/teramono/utilities/pkg/database/models"
	"github.com/teramono/utilities/pkg/messages"
	"github.com/teramono/utilities/pkg/request"
	"github.com/teramono/utilities/pkg/response"
	"gorm.io/gorm"
)

func (server *APIServer) CreateWorkspaces(ctx *gin.Context) {
	// Validate request body.
	body := CreateWorkspacesRequestBody{}
	validErrs, err := server.validateBody(ctx, &body)
	if err != nil {
		response.SetServerErrorResponse(ctx)
		server.logPanicf("%s: %v", messages.InvalidBodyJson, err)
	}

	if len(validErrs) > 0 {
		response.SetUserValidationErrorResponse(ctx, validErrs)
		server.logf("%s - %v", messages.ValidationError.String(), validErrs)
		return
	}

	// Get authorisation credentials from header.
	authCreds := ctx.GetHeader(request.AuthorizationHeader)
	if authCreds == "" {
		response.SetUserErrorResponse(ctx, messages.InvalidAuthorizationHeader)
		server.logf(messages.InvalidAuthorizationHeader.String())
		return
	}

	// Construct message with auth header and workspace name in data.
	data := []byte(fmt.Sprintf(`{"workspaceName":"%s"}`, body.WorkspaceName))
	msg, err := broker.Json(broker.M{(request.AuthorizationHeader): authCreds}, data)
	if err != nil {
		response.SetServerErrorResponse(ctx)
		server.logPanicf("%s: %v", messages.InvalidMessageJson, err)
	}

	// Create a new workspace in the db.
	workspace := models.Workspace{Name: body.WorkspaceName}
	if err := workspace.Create(&server.DB); err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToCreateWorkspace(body.WorkspaceName))
		server.logPanicf("%s: %v", messages.UnableToCreateWorkspace(body.WorkspaceName).String(), err)
	}

	// Publish create subject with workspace id.
	subject := server.subjectOf("create", workspace.ID)
	if err := server.Publish(subject, msg); err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToPublish(subject))
		server.logPanicf("%s: %v", messages.UnableToPublish(subject).String(), err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": messages.Created("workspace"),
	})
}

func (server *APIServer) LoginWorkspaces(ctx *gin.Context) {
	// Get workspace name from header.
	workspaceName := ctx.GetHeader(request.WorkspaceNameHeader)
	if workspaceName == "" {
		response.SetUserErrorResponse(ctx, messages.InvalidWorkspaceNameHeader)
		server.logf(messages.InvalidWorkspaceNameHeader.String())
		return
	}

	// Get authorisation credentials from header.
	authCreds := ctx.GetHeader(request.AuthorizationHeader)
	if authCreds == "" {
		response.SetUserErrorResponse(ctx, messages.InvalidAuthorizationHeader)
		server.logf(messages.InvalidAuthorizationHeader.String())
		return
	}

	// Get workspace by name.
	workspace, err := (&models.Workspace{Name: workspaceName}).GetByName(&server.DB)
	if err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToFindWorkspace(workspaceName))
		server.logPanicf("%s: %v", messages.UnableToFindWorkspace(workspaceName).String(), err)
	}

	// Construct message with auth header.
	msg, err := broker.Json(broker.M{
		(request.AuthorizationHeader): authCreds,
	}, []byte{})
	if err != nil {
		response.SetServerErrorResponse(ctx)
		server.logPanicf("%s: %v", messages.InvalidMessageJson, err)
	}

	// Publish login subject with workspace id.
	subject := server.subjectOf("login", workspace.ID)
	if err := server.Publish(subject, msg); err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToPublish(subject))
		server.logPanicf("%s: %v", messages.UnableToPublish(subject).String(), err)
	}

	ctx.JSON(http.StatusOK, gin.H{ //
		"message": messages.LoggingIn,
	})
}

func (server *APIServer) RunWorkspaces(ctx *gin.Context) {
	// Get workspace id from header.
	workspaceIDStr := ctx.GetHeader(request.WorkspaceIDHeader)
	if workspaceIDStr == "" {
		response.SetUserErrorResponse(ctx, messages.InvalidWorkspaceIDHeader)
		server.logf(messages.InvalidWorkspaceIDHeader.String())
		return
	}

	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		response.SetUserErrorResponse(ctx, messages.InvalidWorkspaceIDHeader)
		server.logf(messages.InvalidWorkspaceIDHeader.String())
		return
	}

	// Find workspace with id in the db.
	exists, err := (&models.Workspace{
		Model: gorm.Model{
			ID: uint(workspaceID),
		},
	}).Exists(&server.DB)
	if err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToFindWorkspace(workspaceIDStr))
		server.logPanicf("%s: %v", messages.UnableToFindWorkspace(workspaceIDStr).String(), err)
	}

	if !exists {
		response.SetUserErrorResponse(ctx, messages.UnableToFindWorkspace(workspaceIDStr))
		server.logf(messages.UnableToFindWorkspace(workspaceIDStr).String())
		return
	}

	// Get all header details.
	header := GetHeadersAsInterface(ctx)

	// Get body details.
	data, err := ctx.GetRawData()
	if err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToGetBody)
		server.logPanicf("%s: %v", messages.UnableToGetBody.String(), err)
	}

	// Construct message with everything from header and body.
	msg, err := broker.Json(header, data)
	if err != nil {
		response.SetServerErrorResponse(ctx)
		server.logPanicf("%s: %v", messages.InvalidMessageJson, err)
	}

	// Publish login subject with workspace id.
	subject := server.subjectOf("run", uint(workspaceID))
	if err := server.Publish(subject, msg); err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToPublish(subject))
		server.logPanicf("%s: %v", messages.UnableToPublish(subject).String(), err)
	}

	// TODO:
	ctx.JSON(http.StatusOK, gin.H{ //
		"message": messages.ScriptStarted,
	})
}
