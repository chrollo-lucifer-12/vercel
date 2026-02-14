package main

import (
	"context"
	"log"

	"github.com/chrollo-lucider-12/proxy/server"
	"github.com/chrollo-lucifer-12/shared/cache"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/storage"
)

func main() {

	env.Load()

	ctx := context.Background()

	dsn := env.Dsn.GetValue()
	db, err := db.NewDB(dsn, ctx)
	if err != nil {
		panic(err)
	}

	cache := cache.NewCacheStore(db)

	st, err := storage.NewS3Storage(env.SupabaseEndpoint.GetValue(), env.SupabaseAccessKey.GetValue(), env.SupabaseAccessSecret.GetValue(), env.Region.GetValue(), "builds")

	s := server.NewServerClient(cache, db, st)

	if err := s.Run(ctx); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}

}
