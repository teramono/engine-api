package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teramono/utilities/pkg/broker"
	"github.com/teramono/utilities/pkg/database/models"
	"github.com/teramono/utilities/pkg/logs"
	"github.com/teramono/utilities/pkg/messages"
	"github.com/teramono/utilities/pkg/request"
	"github.com/teramono/utilities/pkg/response"
)

func (server *APIServer) CreateWorkspaces(ctx *gin.Context) {
	// Validate request body.
	body := CreateWorkspacesRequestBody{}
	validErrs, err := server.validateBody(ctx, &body)
	if err != nil {
		response.SetServerErrorResponse(ctx)
		logs.Panicf(server, "%s: %v", messages.InvalidBodyJson, err)
		return
	}

	if len(validErrs) > 0 {
		response.SetUserValidationErrorResponse(ctx, validErrs)
		logs.Printf(server, "%s - %v", messages.ValidationError.String(), validErrs)
		return
	}

	// Get authorisation credentials from header.
	authCreds := ctx.Request.Header[request.AuthorizationHeader]
	if len(authCreds) < 1 {
		response.SetUserErrorResponse(ctx, messages.InvalidAuthorizationHeader)
		logs.Printf(server, messages.InvalidAuthorizationHeader.String())
		return
	}

	// Construct message with auth header and workspace name in data.
	data := fmt.Sprintf(`{"workspaceName":"%s"}`, body.WorkspaceName)
	msg, err := broker.JsonFromMsgData(
		broker.NewURLFromCtx(ctx),
		broker.Header{
			(request.AuthorizationHeader): authCreds,
		},
		data,
		http.StatusOK,
	)
	if err != nil {
		response.SetServerErrorResponse(ctx)
		logs.Panicf(server, "%s: %v", messages.InvalidMessageJson, err)
		return
	}

	// Create a new workspace in the db.
	workspace := models.Workspace{Name: body.WorkspaceName}
	if err := workspace.Create(&server.DB); err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToCreateWorkspace(body.WorkspaceName))
		logs.Panicf(server, "%s: %v", messages.UnableToCreateWorkspace(body.WorkspaceName).String(), err)
		return
	}

	// Publish create subject with workspace id.
	subject := broker.GetWorkspacesSubjectWithId(&server.Config, "create", workspace.ID)
	if err := server.Publish(subject, msg); err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToPublish(subject))
		logs.Panicf(server, "%s: %v", messages.UnableToPublish(subject).String(), err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": messages.Created("workspace"),
	})
}

func (server *APIServer) Run(ctx *gin.Context) {
	// Get workspace id from header.
	workspaceID, err := server.getWorkspaceID(ctx)
	if err != nil {
		response.SetUserErrorResponse(ctx, messages.InvalidWorkspaceIDHeader)
		logs.Printf(server, "%v", err)
		return
	}

	// Get all header details.
	header := server.getHeaders(ctx)

	// Get body details.
	data, err := ctx.GetRawData()
	if err != nil {
		response.SetServerErrorResponse(ctx, messages.UnableToGetBody)
		logs.Panicf(server, "%s: %v", messages.UnableToGetBody.String(), err)
		return
	}

	// Construct message with everything from header and body.
	msg, err := broker.JsonFromMsgData(broker.NewURLFromCtx(ctx), header, string(data), 0)
	if err != nil {
		response.SetServerErrorResponse(ctx)
		logs.Panicf(server, "%s: %v", messages.InvalidMessageJson, err)
		return
	}

	// Publish login subject with workspace id.
	subject := broker.GetWorkspacesSubjectWithId(&server.Config, "run", workspaceID)
	timeout := time.Duration(server.Config.Engines.API.ReplyTimeout) * time.Second
	reply, err := server.Request(subject, msg, timeout)
	if err != nil {
		response.SetServerErrorResponse(ctx)
		logs.Panicf(server, "%s: %v", messages.UnableToPublish(subject).String(), err)
		return
	}

	// Get response data from server.
	replyMsgData, err := broker.NewMsgData(reply)
	if err != nil {
		response.SetServerErrorResponse(ctx, messages.InvalidReplyJson)
		logs.Panicf(server, "%s: %v", messages.InvalidReplyJson.String(), err)
		return
	}

	// Set headers from reply msg data.
	for k, v := range replyMsgData.Headers {
		ctx.Header(k, strings.Join(v, ";"))
	}

	// TODO: Determine how to returnn value based on Content-Type
	var body gin.H
	if err = json.Unmarshal([]byte(replyMsgData.Data), &body); err != nil {
		response.SetServerErrorResponse(ctx, messages.InvalidReplyDataJson)
		logs.Panicf(server, "%s: %v", messages.InvalidReplyDataJson.String(), err)
		return
	}

	// Set body data and status code from reply msg data.
	ctx.JSON(int(replyMsgData.StatusCode), body)
}
