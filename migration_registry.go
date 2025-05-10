package gormlib

import (
	"sort"
	"sync"
)

// MigrationRegistry gère l'enregistrement et la récupération des migrations
type MigrationRegistry struct {
	migrations map[string]Migration
	mu         sync.RWMutex
}

// NewMigrationRegistry crée un nouveau registre de migrations
func NewMigrationRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		migrations: make(map[string]Migration),
	}
}

// Register enregistre une nouvelle migration
func (r *MigrationRegistry) Register(migration Migration) error {
	if migration == nil {
		return ErrInvalidMigrationName
	}

	name := migration.Name()
	if name == "" {
		return ErrInvalidMigrationName
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.migrations[name]; exists {
		return ErrMigrationAlreadyExists
	}

	r.migrations[name] = migration
	return nil
}

// GetMigrations retourne toutes les migrations enregistrées, triées par nom
func (r *MigrationRegistry) GetMigrations() []Migration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	migrations := make([]Migration, 0, len(r.migrations))
	for _, m := range r.migrations {
		migrations = append(migrations, m)
	}

	// Trier les migrations par nom (qui contient le timestamp)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name() < migrations[j].Name()
	})

	return migrations
}

// GetMigration retourne une migration spécifique par son nom
func (r *MigrationRegistry) GetMigration(name string) (Migration, error) {
	if name == "" {
		return nil, ErrInvalidMigrationName
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if migration, exists := r.migrations[name]; exists {
		return migration, nil
	}
	return nil, ErrMigrationNotFound
}

// HasMigration vérifie si une migration existe
func (r *MigrationRegistry) HasMigration(name string) bool {
	if name == "" {
		return false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.migrations[name]
	return exists
}

// Clear vide le registre de migrations
func (r *MigrationRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.migrations = make(map[string]Migration)
}

// Count retourne le nombre de migrations enregistrées
func (r *MigrationRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.migrations)
} 