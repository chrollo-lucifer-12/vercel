package cache

import (
	"context"
	"fmt"

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
	count := c.db.CheckKey(ctx, key)
	if count != 0 {
		fmt.Println("cache found")
		return nil
	}
	cache := &db.Cache{Key: key, Value: val}
	err := c.db.CreateCache(ctx, cache)
	return err
}

func (c *CacheStore) Get(ctx context.Context, key string) (datatypes.JSON, error) {
	return c.db.GetCache(ctx, key)
}

func (c *CacheStore) Delete(ctx context.Context, key string) error {
	return c.db.DeleteCache(ctx, key)
}
