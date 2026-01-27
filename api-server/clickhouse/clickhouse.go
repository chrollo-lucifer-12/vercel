package clickhouse

import (
	"context"
	"crypto/tls"
	"database/sql"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Log struct {
	Level        string
	Message      string
	CreatedAt    string
	DeploymentID string
}

type View struct {
	DeploymentID string
	Path         string
	ViewDate     time.Time
	Resp         string
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

func (db *ClickHouseDB) GetLogsByDeployment(ctx context.Context, deploymentID string, from, to *time.Time) ([]Log, error) {
	query := `
        SELECT message, level, created_at, deployment_id
        FROM logs
        WHERE deployment_id = ?
    `
	args := []interface{}{deploymentID}

	if from != nil {
		query += " AND created_at >= ?"
		args = append(args, *from)
	}

	// Optional to filter
	if to != nil {
		query += " AND created_at <= ?"
		args = append(args, *to)
	}

	query += `
        ORDER BY created_at DESC
        LIMIT 100
    `

	rows, err := db.conn.QueryContext(ctx, query, args...)
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

func (c *ClickHouseDB) BatchInsertViews(ctx context.Context, views []View) error {
	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO page_view_stats (deployment_id, path, view_date, resp)
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, v := range views {
		if _, err := stmt.ExecContext(ctx, v.DeploymentID, v.Path, v.ViewDate, v.Resp); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (db *ClickHouseDB) GetAnalyticsByDeployment(
	ctx context.Context,
	deploymentID string,
	pagePath string,
	from, to *time.Time,
) ([]View, error) {

	query := `
        SELECT deployment_id, path, view_date, resp
        FROM analytics
        WHERE deployment_id = ?
    `
	args := []interface{}{deploymentID}

	if pagePath != "" {
		query += " AND path = ?"
		args = append(args, pagePath)
	}

	if from != nil {
		query += " AND view_date >= ?"
		args = append(args, *from)
	}

	if to != nil {
		query += " AND view_date <= ?"
		args = append(args, *to)
	}

	query += " ORDER BY view_date DESC LIMIT 100"

	rows, err := db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []View
	for rows.Next() {
		var a View
		if err := rows.Scan(&a.DeploymentID, &a.Path, &a.ViewDate, &a.Resp); err != nil {
			return nil, err
		}
		results = append(results, a)
	}

	return results, nil
}
