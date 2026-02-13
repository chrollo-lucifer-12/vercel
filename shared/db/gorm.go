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

func NewTestDB(dsn string, ctx context.Context) (*DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("No dsn")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
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
	err := d.db.Debug().AutoMigrate(&User{}, &Session{}, &Project{}, &Deployment{}, &LogEvent{}, &Cache{}, &WebsiteAnalytics{})
	if err != nil {
		return err
	}

	_ = d.db.Exec("ALTER TABLE caches SET UNLOGGED").Error

	return nil
}

func first[T any](ctx context.Context, db *gorm.DB, query string, args ...any) (T, error) {
	return gorm.G[T](db).Where(query, args...).First(ctx)
}

func find[T any](ctx context.Context, db *gorm.DB, query string, args ...any) ([]T, error) {
	return gorm.G[T](db).Where(query, args...).Find(ctx)
}

func create[T any](ctx context.Context, db *gorm.DB, val *T) error {
	return gorm.G[T](db).Create(ctx, val)
}

func update[T any](ctx context.Context, db *gorm.DB, query string, val T, args ...any) error {
	_, err := gorm.G[T](db).Where(query, args...).Updates(ctx, val)
	return err
}

func deleteBy[T any](ctx context.Context, db *gorm.DB, query string, args ...any) error {
	_, err := gorm.G[T](db).Where(query, args...).Delete(ctx)
	return err
}

func (d *DB) CreateUser(ctx context.Context, u *User) error {
	return create(ctx, d.db, u)
}

func (d *DB) GetUser(ctx context.Context, email string) (User, error) {
	return first[User](ctx, d.db, "email = ?", email)
}

func (d *DB) UpdateUser(ctx context.Context, u User) error {
	return update[User](ctx, d.db, "id = ?", u, u.ID)
}

func (d *DB) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return deleteBy[User](ctx, d.db, "id = ?", id)
}

func (d *DB) CreateSession(ctx context.Context, s *Session) error {
	return create(ctx, d.db, s)
}

func (d *DB) GetSession(ctx context.Context, id uuid.UUID) (*Session, error) {
	s, err := first[Session](ctx, d.db, "id = ?", id)
	return &s, err
}

func (d *DB) RevokeSession(ctx context.Context, s Session) error {
	return update[Session](ctx, d.db, "id = ?", s, s.ID)
}

func (d *DB) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return deleteBy[Session](ctx, d.db, "id = ?", id)
}

func (d *DB) CreateProject(ctx context.Context, p *Project) error {
	return create(ctx, d.db, p)
}

func (d *DB) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {
	return first[Project](ctx, d.db, "id = ?", id)
}

func (d *DB) GetAllProjects(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	return find[Project](ctx, d.db, "user_id = ?", userID)
}

func (d *DB) UpdateProject(ctx context.Context, id uuid.UUID, p Project) error {
	return update[Project](ctx, d.db, "id = ?", p, id)
}

func (d *DB) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return deleteBy[Project](ctx, d.db, "id = ?", id)
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

func (d *DB) CreateDeployment(ctx context.Context, dep *Deployment) error {
	return create(ctx, d.db, dep)
}

func (d *DB) GetDeploymentByID(ctx context.Context, id uuid.UUID) (Deployment, error) {
	return gorm.G[Deployment](d.db).
		Preload("LogEvents", func(pb gorm.PreloadBuilder) error {
			pb.Order("sequence ASC")
			pb.Limit(500)
			pb.Select("id", "log", "sequence")
			return nil
		}).
		Where("id = ?", id).
		First(ctx)
}

func (d *DB) GetAllDeployments(ctx context.Context, projectID uuid.UUID) ([]Deployment, error) {
	return find[Deployment](ctx, d.db, "project_id = ?", projectID)
}

func (d *DB) UpdateDeployment(ctx context.Context, id uuid.UUID, dep Deployment) error {
	return update[Deployment](ctx, d.db, "id = ?", dep, id)
}

func (d *DB) DeleteDeployment(ctx context.Context, id uuid.UUID) error {
	return deleteBy[Deployment](ctx, d.db, "id = ?", id)
}

func (d *DB) CreateLogEvents(ctx context.Context, logs *[]LogEvent) error {
	return gorm.G[LogEvent](d.db).CreateInBatches(ctx, logs, 10)
}

func (d *DB) GetAnalytics(
	ctx context.Context,
	subdomain string,
	from *time.Time,
	to *time.Time,
) ([]WebsiteAnalytics, error) {

	q := gorm.G[WebsiteAnalytics](d.db).
		Where("subdomain = ?", subdomain)

	if from != nil {
		q = q.Where("created_at >= ?", *from)
	}

	if to != nil {
		q = q.Where("created_at <= ?", *to)
	}

	return q.Find(ctx)
}

func (d *DB) CreateAnalytics(ctx context.Context, w *WebsiteAnalytics) error {
	return create(ctx, d.db, w)
}
