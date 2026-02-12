package main

import (
	"context"

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

	if err := db.MigrateDB(); err != nil {
		panic(err)
	}
}
