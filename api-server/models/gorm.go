package models

import (
	"context"

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

	err = db.AutoMigrate(&Project{}, &Deployment{}, &LogEvent{})
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Raw() *gorm.DB {
	return d.db
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

func (d *DB) CreateLogEvent(ctx context.Context, logEvent *LogEvent) error {
	return gorm.G[LogEvent](d.db).Create(ctx, logEvent)
}

func (d *DB) GetLogEventByID(ctx context.Context, id uuid.UUID) ([]LogEvent, error) {
	logEvent, err := gorm.G[LogEvent](d.db).
		Where("id = ?", id).
		Find(ctx)

	return logEvent, err
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
