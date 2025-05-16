package gormlib

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connection représente une connexion à la base de données
type Connection struct {
	db *gorm.DB
}

// NewConnection crée une nouvelle connexion à la base de données
func NewConnection(config *Config) (*Connection, error) {
	// Configuration du logger GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connexion à la base de données
	db, err := gorm.Open(postgres.Open(config.DSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la connexion à la base de données: %v", err)
	}

	// Configuration du pool de connexions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la connexion SQL: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return &Connection{db: db}, nil
}

// DB retourne l'instance de GORM
func (c *Connection) DB() *gorm.DB {
	return c.db
}

// Close ferme la connexion à la base de données
func (c *Connection) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de la connexion SQL: %v", err)
	}
	return sqlDB.Close()
}
