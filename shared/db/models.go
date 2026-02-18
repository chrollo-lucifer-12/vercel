package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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

type User struct {
	Base
	Name       string    `gorm:"not null;check:name <> ''"`
	Email      string    `gorm:"unique;not null;check:email <> ''"`
	Password   string    `gorm:"not null;check:password <> ''"`
	Projects   []Project `gorm:"foreignKey:UserID"`
	IsVerified bool
}

type Otp struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	User      *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Token     string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time
}

type Session struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	User         *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	UserEmail    string    `gorm:"not null"`
	RefreshToken string    `gorm:"not null"`
	Revoked      bool
	ExpiresAt    time.Time
}

type Project struct {
	Base
	Name         string
	GitUrl       string
	SubDomain    string
	CustomDomain string
	UserID       uuid.UUID
	Deployments  []Deployment `gorm:"foreignKey:ProjectID"`
}

type Deployment struct {
	Base
	ProjectID uuid.UUID `gorm:"type:uuid;index"`
	Status    string
	LogEvents []LogEvent `gorm:"foreignKey:DeploymentID"`
	Sequence  int        `gorm:"autoIncrement"`
}

type LogEvent struct {
	Base
	DeploymentID uuid.UUID `gorm:"type:uuid;index"`
	Log          string
	Metadata     datatypes.JSON `gorm:"type:jsonb"`
	Sequence     int64
}

type WebsiteAnalytics struct {
	Base
	Subdomain      string `gorm:"type:varchar(255);not null;index:idx_subdomain" json:"subdomain"`
	Path           string `gorm:"type:varchar(1000);not null" json:"path"`
	Method         string `gorm:"type:varchar(10);not null" json:"method"`
	StatusCode     int    `gorm:"not null" json:"status_code"`
	ResponseTimeMs int    `gorm:"type:int" json:"response_time_ms,omitempty"`
	UserAgent      string `gorm:"type:text" json:"user_agent"`
	IPAddress      string `gorm:"type:varchar(45)" json:"ip_address"`
	Referer        string `gorm:"type:text" json:"referer"`
}

type Cache struct {
	Base
	Key   string         `gorm:"unique;index"`
	Value datatypes.JSON `gorm:"type:jsonb"`
}
