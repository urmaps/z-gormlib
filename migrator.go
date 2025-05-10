package gormlib

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// Migrator gère les migrations de la base de données
type Migrator struct {
	db     *gorm.DB
	config *MigrationConfig
}

// NewMigrator crée un nouveau gestionnaire de migrations
func NewMigrator(db *gorm.DB, config *MigrationConfig) *Migrator {
	if config == nil {
		config = DefaultConfig()
	}
	return &Migrator{
		db:     db,
		config: config,
	}
}

// RunMigrations exécute toutes les migrations non appliquées
func (m *Migrator) RunMigrations(migrations ...Migration) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	// Créer la table des migrations si elle n'existe pas
	if err := m.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return NewMigrationError("create migrations table", err)
	}

	// Exécuter les migrations par lots
	for i := 0; i < len(migrations); i += m.config.BatchSize {
		end := i + m.config.BatchSize
		if end > len(migrations) {
			end = len(migrations)
		}

		batch := migrations[i:end]
		if err := m.runMigrationBatch(ctx, batch); err != nil {
			return err
		}
	}

	return nil
}

// runMigrationBatch exécute un lot de migrations dans une transaction
func (m *Migrator) runMigrationBatch(ctx context.Context, migrations []Migration) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, migration := range migrations {
			// Vérifier si la migration a déjà été appliquée
			var record MigrationRecord
			result := tx.Where("name = ?", migration.Name()).First(&record)
			if result.Error == nil {
				continue
			}

			// Exécuter la migration avec retry
			var err error
			for attempt := 0; attempt < m.config.RetryAttempts; attempt++ {
				if err = migration.Up(tx); err == nil {
					break
				}
				time.Sleep(time.Second * time.Duration(attempt+1))
			}
			if err != nil {
				return NewMigrationError("run migration", err)
			}

			// Enregistrer la migration
			record = MigrationRecord{
				Name:      migration.Name(),
				AppliedAt: time.Now(),
			}
			if err := tx.Create(&record).Error; err != nil {
				return NewMigrationError("record migration", err)
			}
		}
		return nil
	})
}

// RollbackMigration annule la dernière migration
func (m *Migrator) RollbackMigration(migration Migration) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
	defer cancel()

	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Vérifier si la migration existe
		var record MigrationRecord
		if err := tx.Where("name = ?", migration.Name()).First(&record).Error; err != nil {
			return ErrMigrationNotFound
		}

		// Exécuter le rollback avec retry
		var err error
		for attempt := 0; attempt < m.config.RetryAttempts; attempt++ {
			if err = migration.Down(tx); err == nil {
				break
			}
			time.Sleep(time.Second * time.Duration(attempt+1))
		}
		if err != nil {
			return NewMigrationError("rollback migration", err)
		}

		// Supprimer l'enregistrement de la migration
		if err := tx.Delete(&record).Error; err != nil {
			return NewMigrationError("delete migration record", err)
		}

		return nil
	})
}

// GetAppliedMigrations retourne la liste des migrations appliquées
func (m *Migrator) GetAppliedMigrations() ([]MigrationRecord, error) {
	var migrations []MigrationRecord
	if err := m.db.Order("applied_at").Find(&migrations).Error; err != nil {
		return nil, NewMigrationError("get applied migrations", err)
	}
	return migrations, nil
}

// GetPendingMigrations retourne la liste des migrations en attente
func (m *Migrator) GetPendingMigrations(availableMigrations []Migration) ([]Migration, error) {
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]bool)
	for _, m := range applied {
		appliedMap[m.Name] = true
	}

	var pending []Migration
	for _, m := range availableMigrations {
		if !appliedMap[m.Name()] {
			pending = append(pending, m)
		}
	}

	return pending, nil
} 