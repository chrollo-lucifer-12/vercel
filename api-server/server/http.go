package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/workflow"
)

type ServerClient struct {
	wClient *workflow.WorkflowClient
	db      *db.DB
}

func NewServerClient(wClient *workflow.WorkflowClient, db *db.DB) (*ServerClient, error) {
	if wClient == nil {
		return nil, fmt.Errorf("No workflow client")
	}
	if db == nil {
		return nil, fmt.Errorf("No db client")
	}
	return &ServerClient{wClient: wClient, db: db}, nil
}

func (h *ServerClient) StartHTTP() {

	http.HandleFunc("/api/v1/deploy", h.deployHandler)
	http.HandleFunc("/api/v1/project", h.projectHandler)

	log.Fatal(http.ListenAndServe(":9000", nil))
}
