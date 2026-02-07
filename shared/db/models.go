package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt time.Time `gorm:"index"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type Project struct {
	Base
	Name         string
	GitUrl       string
	SubDomain    string
	CustomDomain string
	Deployments  []Deployment `gorm:"foreignKey:ProjectID"`
}

type Deployment struct {
	Base
	ProjectID uuid.UUID `gorm:"type:uuid;index"`
	Status    string
	LogEvents []LogEvent `gorm:"foreignKey:DeploymentID"`
}

type GitHash struct {
	Base
	ProjectID uuid.UUID `gorm:"type:uuid;index"`
	Hash      string
}

type LogEvent struct {
	Base
	DeploymentID uuid.UUID `gorm:"type:uuid;index"`
	Log          string
}

type WebsiteAnalytics struct {
	Base
	Status       int
	Method       string
	OriginalPath string
	Slug         string
}
