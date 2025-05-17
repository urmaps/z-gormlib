package gormlib

import (
	"time"

	"gorm.io/gorm"
)

// Migration représente une migration de base de données
type Migration interface {
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
	Name() string
}

// MigrationRecord représente une migration appliquée dans la base de données
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	AppliedAt time.Time `gorm:"not null"`
}
