package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type DB struct {
	db *gorm.DB
}

func NewDB(dsn string, ctx context.Context) (*DB, error) {

	if dsn == "" {
		return nil, fmt.Errorf("No dsn")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Raw() *gorm.DB {
	return d.db
}

func (d *DB) MigrateDB() error {
	err := d.db.AutoMigrate(&Project{}, &Deployment{}, &LogEvent{}, &GitHash{}, &Cache{}, &User{})
	if err != nil {
		return err
	}

	err = d.db.Exec("ALTER TABLE caches SET UNLOGGED").Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetAllDeployments(ctx context.Context, projectID uuid.UUID) ([]Deployment, error) {
	return gorm.G[Deployment](d.db).Where("project_id = ?", projectID).Find(ctx)
}

func (d *DB) GetAllProjects(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	return gorm.G[Project](d.db).Where("user_id = ?", userID).Find(ctx)
}

func (d *DB) CreateUser(ctx context.Context, u *User) error {
	err := gorm.G[User](d.db).Create(ctx, u)
	return err
}

func (d *DB) GetUser(ctx context.Context, email string) (User, error) {
	return gorm.G[User](d.db).Where("email = ?", email).First(ctx)
}

func (d *DB) UpdateUser(ctx context.Context, u User) error {
	_, err := gorm.G[User](d.db).Updates(ctx, u)

	return err
}

func (d *DB) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[User](d.db).Where("id = ?", id).Delete(ctx)
	return err
}

func (d *DB) CreateSession(ctx context.Context, s *Session) error {
	return gorm.G[Session](d.db).Create(ctx, s)
}

func (d *DB) GetSession(ctx context.Context, id uuid.UUID) (*Session, error) {
	session, err := gorm.G[Session](d.db).Where("id = ?", id).First(ctx)
	return &session, err
}

func (d *DB) RevokeSession(ctx context.Context, s Session) error {
	_, err := gorm.G[Session](d.db).Updates(ctx, s)
	return err
}

func (d *DB) DeleteSession(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[Session](d.db).Where("id = ?", id).Delete(ctx)
	return err
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

func (d *DB) CreateCache(ctx context.Context, cache *Cache) error {
	return d.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(cache).Error
}

func (d *DB) CheckKey(ctx context.Context, key string) int64 {
	count, _ := gorm.G[Cache](d.db).Where("key = ?", key).Count(ctx, key)
	return count
}

func (d *DB) DeleteCache(ctx context.Context, key string) error {
	_, err := gorm.G[Cache](d.db).Where("key = ?", key).Delete(ctx)
	return err
}

func (d *DB) GetCache(ctx context.Context, key string) (datatypes.JSON, error) {
	cache, err := gorm.G[Cache](d.db).Where("key = ?", key).First(ctx)
	return cache.Value, err
}

func (d *DB) GetDeploymentByID(ctx context.Context, id uuid.UUID) (Deployment, error) {
	deployment, err := gorm.G[Deployment](d.db).Preload("LogEvents", func(pb gorm.PreloadBuilder) error {
		pb.Order("sequence ASC")
		pb.Limit(500)
		pb.Select("id", "log", "sequence")
		return nil
	}).
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

func (d *DB) CreateAnalytcis(ctx context.Context, w *WebsiteAnalytics) error {
	return gorm.G[WebsiteAnalytics](d.db).Create(ctx, w)
}
