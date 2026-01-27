package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrollo-lucifer-12/api-server/clickhouse"
	"github.com/chrollo-lucifer-12/api-server/models"
	"github.com/chrollo-lucifer-12/api-server/workflow"
)

type ServerClient struct {
	wClient *workflow.WorkflowClient
	db      *models.DB
	clickDB *clickhouse.ClickHouseDB
}

func NewServerClient(wClient *workflow.WorkflowClient, db *models.DB, clickDB *clickhouse.ClickHouseDB) (*ServerClient, error) {
	if wClient == nil {
		return nil, fmt.Errorf("No workflow client")
	}
	return &ServerClient{wClient: wClient, db: db, clickDB: clickDB}, nil
}

func (h *ServerClient) StartHTTP() {
	http.HandleFunc("/deploy", h.deployHandler)
	http.HandleFunc("/project", h.projectHandler)
	http.HandleFunc("/logs/", h.logsHandler)
	http.HandleFunc("/analytics/", h.analyticsHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
