package gormlib

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config contient la configuration de la base de données
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	Schema          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// NewConfig crée une nouvelle configuration à partir des variables d'environnement
func NewConfig() *Config {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))
	connMaxLifetime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "1h"))

	return &Config{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            port,
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres"),
		Database:        getEnv("DB_NAME", "oauth"),
		Schema:          getEnv("DB_SCHEMA", "public"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxIdleConns:    maxIdleConns,
		MaxOpenConns:    maxOpenConns,
		ConnMaxLifetime: connMaxLifetime,
	}
}

// DSN retourne la chaîne de connexion PostgreSQL
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.Schema, c.SSLMode)
}

// getEnv récupère une variable d'environnement avec une valeur par défaut
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// MigrationConfig contient la configuration pour les migrations
type MigrationConfig struct {
	// BatchSize est le nombre de migrations à exécuter en une seule transaction
	BatchSize int

	// Timeout est le délai maximum pour l'exécution d'une migration
	Timeout time.Duration

	// RetryAttempts est le nombre de tentatives en cas d'échec
	RetryAttempts int

	// TableName est le nom de la table qui stocke les migrations
	TableName string

	// AutoCreateDir indique si le dossier des migrations doit être créé automatiquement
	AutoCreateDir bool
}

// DefaultConfig retourne la configuration par défaut
func DefaultConfig() *MigrationConfig {
	return &MigrationConfig{
		BatchSize:     10,
		Timeout:       5 * time.Minute,
		RetryAttempts: 3,
		TableName:     "migrations",
		AutoCreateDir: true,
	}
}

// Constants
const (
	// MigrationFileSuffix est le suffixe des fichiers de migration
	MigrationFileSuffix = ".go"

	// MigrationStructPrefix est le préfixe des structures de migration
	MigrationStructPrefix = "Migration"

	// DefaultMigrationsDir est le dossier par défaut pour les migrations
	DefaultMigrationsDir = "migrations"

	// DefaultFileMode est le mode par défaut pour les fichiers de migration
	DefaultFileMode = 0644

	// DefaultDirMode est le mode par défaut pour les dossiers de migration
	DefaultDirMode = 0755
)
