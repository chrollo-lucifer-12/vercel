package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	db *gorm.DB
}

func NewDB(dsn string, ctx context.Context) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Project{}, &Deployment{}, &LogEvent{}, &GitHash{})
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Raw() *gorm.DB {
	return d.db
}

func (d *DB) CreateHash(ctx context.Context, gitHash *GitHash) error {
	return gorm.G[GitHash](d.db).Create(ctx, gitHash)
}

func (d *DB) UpdateHash(ctx context.Context, project_id uuid.UUID, gitHash GitHash) error {
	_, err := gorm.G[GitHash](d.db).Where("project_id = ?", project_id).Updates(ctx, gitHash)
	return err
}

func (d *DB) CreateProject(ctx context.Context, project *Project) error {
	return gorm.G[Project](d.db).Create(ctx, project)
}

func (d *DB) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {

	p, err := gorm.G[Project](d.db).
		Where("id = ?", id).
		First(ctx)

	return p, err
}

func (d *DB) UpdateProject(
	ctx context.Context,
	id uuid.UUID,
	p Project,
) error {
	_, err := gorm.G[Project](d.db).
		Where("id = ?", id).
		Updates(ctx, p)
	return err
}

func (d *DB) DeleteProject(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[Project](d.db).
		Where("id = ?", id).
		Delete(ctx)

	return err
}

func (d *DB) CreateDeployment(ctx context.Context, deployment *Deployment) error {
	return gorm.G[Deployment](d.db).Create(ctx, deployment)
}

func (d *DB) GetDeploymentByID(ctx context.Context, id uuid.UUID) (Deployment, error) {
	deployment, err := gorm.G[Deployment](d.db).
		Where("id = ?", id).
		First(ctx)

	return deployment, err
}

func (d *DB) UpdateDeployment(
	ctx context.Context,
	id uuid.UUID,
	deployment Deployment,
) error {
	_, err := gorm.G[Deployment](d.db).
		Where("id = ?", id).
		Updates(ctx, deployment)

	return err
}

func (d *DB) DeleteDeployment(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[Deployment](d.db).
		Where("id = ?", id).
		Delete(ctx)

	return err
}

func (d *DB) CreateLogEvents(ctx context.Context, logEvents *[]LogEvent) error {
	return gorm.G[LogEvent](d.db).CreateInBatches(ctx, logEvents, 10)
}

func (d *DB) GetLogEventByID(ctx context.Context, id uuid.UUID) ([]LogEvent, error) {
	logEvent, err := gorm.G[LogEvent](d.db).
		Where("deployment_id = ?", id).
		Find(ctx)

	return logEvent, err
}

func (d *DB) GetLogEventsByDeploymentAndTimeRange(
	ctx context.Context,
	deploymentID uuid.UUID,
	from time.Time,
	to time.Time,
) ([]LogEvent, error) {

	logEvents, err := gorm.G[LogEvent](d.db).
		Where(
			"message!=analytics AND deployment_id = ? AND created_at BETWEEN ? AND ?",
			deploymentID,
			from,
			to,
		).
		Find(ctx)

	return logEvents, err
}

func (d *DB) GetAnalytics(ctx context.Context, deploymentID uuid.UUID, from time.Time, to time.Time, status_code string, path string) ([]LogEvent, error) {
	q := gorm.G[LogEvent](d.db).
		Where("deployment_id = ?", deploymentID).
		Where("created_at BETWEEN ? AND ?", from, to)

	if status_code != "" {
		q = q.Where(
			"metadata ->> 'status_code' = ?",
			status_code,
		)
	}

	if path != "" {
		q = q.Where(
			"metadata ->> 'path' = ?",
			path,
		)
	}

	logs, err := q.Find(ctx)
	return logs, err
}

func (d *DB) UpdateLogEvent(
	ctx context.Context,
	id uuid.UUID,
	logEvent LogEvent,
) error {
	_, err := gorm.G[LogEvent](d.db).
		Where("id = ?", id).
		Updates(ctx, logEvent)

	return err
}

func (d *DB) DeleteLogEvent(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[LogEvent](d.db).
		Where("id = ?", id).
		Delete(ctx)

	return err
}
