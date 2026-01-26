package models

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB(dsn string, ctx context.Context) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Project{}, &Deployment{}, &LogEvent{})
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}
