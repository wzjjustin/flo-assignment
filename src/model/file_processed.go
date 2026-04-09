package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProcessedFile struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;primaryKey"`
	Checksum string    `gorm:"type:varchar(64);not null;unique"`
}
