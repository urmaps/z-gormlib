package gormlib

import (
	"fmt"
	"sort"
	"sync"
)

// MigrationRegistry gère l'enregistrement et le suivi des migrations
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

// globalRegistry est le registre global des migrations
var globalRegistry = NewMigrationRegistry()

// Register enregistre une nouvelle migration dans le registre
func (r *MigrationRegistry) Register(migration Migration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := migration.Name()
	if name == "" {
		return fmt.Errorf("le nom de la migration ne peut pas être vide")
	}

	if _, exists := r.migrations[name]; exists {
		return fmt.Errorf("une migration avec le nom %s existe déjà", name)
	}

	r.migrations[name] = migration
	return nil
}

// GetMigrationByName retourne une migration par son nom
func (r *MigrationRegistry) GetMigrationByName(name string) Migration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.migrations[name]
}

// GetAllMigrations retourne toutes les migrations enregistrées, triées par nom
func (r *MigrationRegistry) GetAllMigrations() []Migration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	migrations := make([]Migration, 0, len(r.migrations))
	for _, m := range r.migrations {
		migrations = append(migrations, m)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name() < migrations[j].Name()
	})

	return migrations
}

// RegisterGlobal enregistre une migration dans le registre global
func RegisterGlobal(migration Migration) error {
	return globalRegistry.Register(migration)
}

// GetGlobalMigrationByName retourne une migration du registre global par son nom
func GetGlobalMigrationByName(name string) Migration {
	return globalRegistry.GetMigrationByName(name)
}

// GetAllGlobalMigrations retourne toutes les migrations du registre global
func GetAllGlobalMigrations() []Migration {
	return globalRegistry.GetAllMigrations()
}

// HasMigration vérifie si une migration existe dans le registre
func (r *MigrationRegistry) HasMigration(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.migrations[name]
	return exists
}

// Clear vide le registre de migrations (utile pour les tests)
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
