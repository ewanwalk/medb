package models

import (
	"github.com/Ewan-Walker/gorm"
	"time"
)

type Setting struct {
	ID      int64     `gorm:"AUTO_INCREMENT;primary_key"`
	Name    string    `gorm:"type:varchar(255);not null"`
	Value   string    `gorm:"type:varchar(255);not null"`
	Created time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Updated time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (s *Setting) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("Updated", time.Now().UTC())
}
