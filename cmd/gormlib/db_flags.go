package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/urmaps/z-gormlib"
)

// DBFlags représente les flags liés à la configuration de la base de données
type DBFlags struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLifetime time.Duration
}

// ParseDBFlags parse les flags de la ligne de commande pour la configuration de la base de données
func ParseDBFlags() *DBFlags {
	dbFlags := &DBFlags{}

	// Flags de base de données
	flag.StringVar(&dbFlags.Host, "db-host", getEnv("DB_HOST", "localhost"), "Database host")
	flag.IntVar(&dbFlags.Port, "db-port", getEnvAsInt("DB_PORT", 5432), "Database port")
	flag.StringVar(&dbFlags.User, "db-user", getEnv("DB_USER", "postgres"), "Database user")
	flag.StringVar(&dbFlags.Password, "db-password", getEnv("DB_PASSWORD", "postgres"), "Database password")
	flag.StringVar(&dbFlags.Database, "db-name", getEnv("DB_NAME", "postgres"), "Database name")
	flag.StringVar(&dbFlags.SSLMode, "db-ssl-mode", getEnv("DB_SSL_MODE", "disable"), "Database SSL mode")
	flag.IntVar(&dbFlags.MaxOpenConns, "db-max-open-conns", getEnvAsInt("DB_MAX_OPEN_CONNS", 25), "Maximum number of open connections")
	flag.IntVar(&dbFlags.MaxIdleConns, "db-max-idle-conns", getEnvAsInt("DB_MAX_IDLE_CONNS", 25), "Maximum number of idle connections")
	flag.DurationVar(&dbFlags.ConnMaxLifetime, "db-conn-max-lifetime", getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute), "Maximum connection lifetime")

	return dbFlags
}

// ToConfig convertit les flags en configuration GORM
func (f *DBFlags) ToConfig() *gormlib.Config {
	config := gormlib.DefaultConfig()
	config.SetDSN(f.buildDSN())
	config.SetMaxOpenConns(f.MaxOpenConns)
	config.SetMaxIdleConns(f.MaxIdleConns)
	config.SetConnMaxLifetime(f.ConnMaxLifetime)
	return config
}

// buildDSN construit la chaîne de connexion DSN
func (f *DBFlags) buildDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		f.Host, f.Port, f.User, f.Password, f.Database, f.SSLMode)
}

// getEnv récupère une variable d'environnement avec une valeur par défaut
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt récupère une variable d'environnement comme un entier
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration récupère une variable d'environnement comme une durée
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
} 