package cache

import (
	"context"

	"github.com/chrollo-lucifer-12/shared/db"
	"gorm.io/datatypes"
)

type MockCacheDB struct {
	CreateCacheFn func(ctx context.Context, cache *db.Cache) error
	GetCacheFn    func(ctx context.Context, key string) (datatypes.JSON, error)
	DeleteCacheFn func(ctx context.Context, key string) error
}

func (m *MockCacheDB) CreateCache(ctx context.Context, cache *db.Cache) error {
	return m.CreateCacheFn(ctx, cache)
}

func (m *MockCacheDB) GetCache(ctx context.Context, key string) (datatypes.JSON, error) {
	return m.GetCacheFn(ctx, key)
}

func (m *MockCacheDB) DeleteCache(ctx context.Context, key string) error {
	return m.DeleteCacheFn(ctx, key)
}

type CacheDB interface {
	CreateCache(ctx context.Context, cache *db.Cache) error
	GetCache(ctx context.Context, key string) (datatypes.JSON, error)
	DeleteCache(ctx context.Context, key string) error
}

type CacheStore struct {
	db CacheDB
}

func NewCacheStore(db CacheDB) *CacheStore {
	return &CacheStore{db: db}
}

func (c *CacheStore) Set(ctx context.Context, key string, val datatypes.JSON) error {

	cache := &db.Cache{
		Key:   key,
		Value: val,
	}

	return c.db.CreateCache(ctx, cache)
}

func (c *CacheStore) Get(ctx context.Context, key string) (datatypes.JSON, error) {
	return c.db.GetCache(ctx, key)
}

func (c *CacheStore) Delete(ctx context.Context, key string) error {
	return c.db.DeleteCache(ctx, key)
}
