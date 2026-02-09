package cache

import (
	"context"

	"github.com/chrollo-lucifer-12/shared/db"
	"gorm.io/datatypes"
)

type CacheStore struct {
	db *db.DB
}

func NewCacheStore(db *db.DB) *CacheStore {
	return &CacheStore{db: db}
}

func (c *CacheStore) Set(ctx context.Context, key string, val datatypes.JSON) error {
	cache := &db.Cache{Key: key, Value: val}
	err := c.db.CreateCache(ctx, cache)
	return err
}

func (c *CacheStore) Get(ctx context.Context, key string) (datatypes.JSON, error) {
	return c.db.GetCache(ctx, key)
}
