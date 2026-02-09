package main

import (
	"context"
	"log"

	"github.com/chrollo-lucider-12/proxy/server"
	"github.com/chrollo-lucifer-12/shared/cache"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/chrollo-lucifer-12/shared/env"
)

func main() {

	err := env.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	dsn := env.Dsn.GetValue()
	db, err := db.NewDB(dsn, ctx)
	if err != nil {
		panic(err)
	}

	cache := cache.NewCacheStore(db)

	s := server.NewServerClient(cache)

	if err := s.Run(ctx); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}

}
