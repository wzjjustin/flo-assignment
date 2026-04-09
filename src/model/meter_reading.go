package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MeterReading struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;primaryKey"`
	NMI         string    `gorm:"type:varchar(10);not null;uniqueIndex:idx_nmi_timestamp"`
	Timestamp   time.Time `gorm:"type:timestamp;not null;uniqueIndex:idx_nmi_timestamp"`
	Consumption string    `gorm:"type:numeric;not null"`
}
