package models

import (
	"time"
)

type Setting struct {
	ID        int64     `gorm:"AUTO_INCREMENT;primary_key"`
	Name      string    `gorm:"type:varchar(255);not null;unique_index"`
	Value     string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
