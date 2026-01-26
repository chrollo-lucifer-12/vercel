package clickhouse

import (
	"context"
	"crypto/tls"
	"database/sql"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Log struct {
	Level        string
	Message      string
	CreatedAt    string
	DeploymentID string
}

type ClickHouseDB struct {
	conn *sql.DB
}

func NewClickHouseDB(addr, username, password string) *ClickHouseDB {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr:     []string{addr},
		Protocol: clickhouse.Native,
		TLS:      &tls.Config{},
		Auth: clickhouse.Auth{
			Username: username,
			Password: password,
		},
	})
	return &ClickHouseDB{conn: conn}
}

func (c *ClickHouseDB) BatchInsertLogs(ctx context.Context, logs []Log) error {
	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO logs (deployment_id, message, level, created_at)
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, l := range logs {
		if _, err := stmt.ExecContext(ctx, l.DeploymentID, l.Message, l.Level, l.CreatedAt); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (db *ClickHouseDB) GetLogsByDeployment(ctx context.Context, deploymentID string) ([]Log, error) {
	rows, err := db.conn.QueryContext(ctx, `
        SELECT  message, level, created_at, deployment_id
        FROM logs
        WHERE deployment_id = ?
        ORDER BY created_at DESC
        LIMIT 100
    `, deploymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []Log

	for rows.Next() {
		var l Log
		if err := rows.Scan(&l.Message, &l.Level, &l.CreatedAt, &l.DeploymentID); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	return logs, nil
}
