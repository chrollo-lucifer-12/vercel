package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrollo-lucifer-12/api-server/models"
	"github.com/chrollo-lucifer-12/api-server/workflow"
	"github.com/chrollo-lucifer-12/api-server/ws"
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
	http.HandleFunc("/ws", ws.WsHandler)
	http.HandleFunc("/deploy", h.deployHandler)
	http.HandleFunc("/project", h.projectHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
