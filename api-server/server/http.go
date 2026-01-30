package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrollo-lucifer-12/api-server/models"
	"github.com/chrollo-lucifer-12/api-server/workflow"
)

type ServerClient struct {
	wClient *workflow.WorkflowClient
	db      *models.DB
}

func NewServerClient(wClient *workflow.WorkflowClient, db *models.DB) (*ServerClient, error) {
	if wClient == nil {
		return nil, fmt.Errorf("No workflow client")
	}
	return &ServerClient{wClient: wClient, db: db}, nil
}

func (h *ServerClient) StartHTTP() {

	http.HandleFunc("/api/v1/deploy", h.deployHandler)
	http.HandleFunc("/api/v1/project", h.projectHandler)
	http.HandleFunc("/api/v1/logs", h.logsHandler)
	http.HandleFunc("/api/v1/logs/insert", h.registerLogsRoutes)
	http.HandleFunc("/api/v1/analytics", h.analyticsHandler)

	log.Fatal(http.ListenAndServe(":9000", nil))
}
