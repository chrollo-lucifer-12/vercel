package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type User struct {
	Base
	Name       string    `gorm:"not null;check:name <> ''" json:"name"`
	Email      string    `gorm:"unique;not null;check:email <> ''" json:"email"`
	Password   string    `gorm:"not null;check:password <> ''" json:"-"`
	Projects   []Project `gorm:"foreignKey:UserID" json:"projects,omitempty"`
	IsVerified bool      `json:"is_verified"`
}

type Otp struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Token     string    `json:"token"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Session struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	User         *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	UserEmail    string    `gorm:"not null" json:"user_email"`
	RefreshToken string    `gorm:"not null" json:"refresh_token"`
	Revoked      bool      `json:"revoked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Project struct {
	Base
	Name         string       `json:"name"`
	GitUrl       string       `json:"git_url"`
	SubDomain    string       `json:"sub_domain"`
	CustomDomain string       `json:"custom_domain"`
	UserID       uuid.UUID    `json:"user_id"`
	Deployments  []Deployment `gorm:"foreignKey:ProjectID" json:"deployments,omitempty"`
}

type Deployment struct {
	Base
	ProjectID uuid.UUID  `gorm:"type:uuid;index" json:"project_id"`
	Status    string     `json:"status"`
	LogEvents []LogEvent `gorm:"foreignKey:DeploymentID" json:"log_events,omitempty"`
	Sequence  int        `gorm:"autoIncrement" json:"sequence"`
}

type LogEvent struct {
	Base
	DeploymentID uuid.UUID      `gorm:"type:uuid;index;not null" json:"deployment_id"`
	Deployment   Deployment     `gorm:"foreignKey:DeploymentID" json:"deployment,omitempty"`
	Log          string         `json:"log"`
	Metadata     datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	Sequence     int64          `json:"sequence"`
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
