package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/teramono/utilities/pkg/broker"
)

func (server *APIServer) subjectOf(action string, id uint) string {
	version := server.Config.Broker.Subscriptions.Workspaces.Version
	return broker.SubjectOf(version, action, "workspaces", id)
}

func (server *APIServer) logPanicf(format string, v ...interface{}) {
	version := server.Config.Broker.Subscriptions.Logs.Version
	subject := broker.SubjectOf(version, "create", "logs", "panic")
	message := fmt.Sprintf(format, v...)

	if err := server.Publish(subject, []byte(message)); err != nil {
		log.Printf("unable to publish error log: %v", err)
	}

	log.Panicf(format, v...)
}

func (server *APIServer) logf(format string, v ...interface{}) {
	version := server.Config.Broker.Subscriptions.Logs.Version
	subject := broker.SubjectOf(version, "create", "logs", "panic")
	message := fmt.Sprintf(format, v...)

	if err := server.Publish(subject, []byte(message)); err != nil {
		log.Printf("unable to publish error log: %v", err)
	}

	log.Printf(format, v...)
}

func GetHeadersAsInterface(ctx *gin.Context) map[string]interface{} {
	headers := map[string]interface{}{}

	for k, v := range ctx.Request.Header {
		headers[k] = v
	}

	return headers
}
