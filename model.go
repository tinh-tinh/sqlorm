package sqlorm

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// deprecated
type Model struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt *time.Time     `gorm:"not null;default:now()"`
	UpdatedAt *time.Time     `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// deprecated
type ModelSnakeCase struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	CreatedAt *time.Time     `gorm:"not null;default:now()" json:"created_at,omitempty"`
	UpdatedAt *time.Time     `gorm:"not null;default:now()" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// deprecated
type ModelCamelCase struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	CreatedAt *time.Time     `gorm:"not null;default:now()" json:"createdAt,omitempty"`
	UpdatedAt *time.Time     `gorm:"not null;default:now()" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
