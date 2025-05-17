package gormlib

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MigrationDiscovery gère la découverte automatique des migrations
type MigrationDiscovery struct {
	MigrationsDir string
	registry      *MigrationRegistry
}

// NewMigrationDiscovery crée un nouveau découvreur de migrations
func NewMigrationDiscovery(migrationsDir string) *MigrationDiscovery {
	return &MigrationDiscovery{
		MigrationsDir: migrationsDir,
		registry:      globalRegistry,
	}
}

// DiscoverMigrations découvre et charge automatiquement toutes les migrations
// dans le dossier spécifié, triées par ordre chronologique
func (d *MigrationDiscovery) DiscoverMigrations() ([]Migration, error) {
	// Créer le dossier migrations s'il n'existe pas
	if err := os.MkdirAll(d.MigrationsDir, 0755); err != nil {
		return nil, fmt.Errorf("erreur lors de la création du dossier migrations: %v", err)
	}

	// Lire tous les fichiers .go dans le dossier migrations
	files, err := os.ReadDir(d.MigrationsDir)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du dossier migrations: %v", err)
	}

	type migrationInfo struct {
		migration Migration
		timestamp time.Time
	}

	var migrationsInfo []migrationInfo
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			// Extraire le timestamp et le nom de la migration
			timestamp, name, err := d.parseMigrationFileName(file.Name())
			if err != nil {
				continue // Skip les fichiers qui ne suivent pas le format
			}

			// Vérifier si la migration est déjà enregistrée
			migration := d.registry.GetMigrationByName(name)
			if migration != nil {
				migrationsInfo = append(migrationsInfo, migrationInfo{
					migration: migration,
					timestamp: timestamp,
				})
			}
		}
	}

	// Trier les migrations par timestamp
	sort.Slice(migrationsInfo, func(i, j int) bool {
		return migrationsInfo[i].timestamp.Before(migrationsInfo[j].timestamp)
	})

	// Convertir en slice de Migration
	migrations := make([]Migration, len(migrationsInfo))
	for i, info := range migrationsInfo {
		migrations[i] = info.migration
	}

	return migrations, nil
}

// parseMigrationFileName extrait le timestamp et le nom de la migration du nom de fichier
func (d *MigrationDiscovery) parseMigrationFileName(filename string) (time.Time, string, error) {
	// Format attendu: YYYYMMDDHHMMSS_name.go
	base := strings.TrimSuffix(filename, ".go")
	parts := strings.SplitN(base, "_", 2)
	if len(parts) != 2 {
		return time.Time{}, "", fmt.Errorf("format de nom de fichier invalide: %s", filename)
	}

	timestamp, err := time.Parse("20060102150405", parts[0])
	if err != nil {
		return time.Time{}, "", fmt.Errorf("timestamp invalide dans le nom de fichier: %s", filename)
	}

	return timestamp, base, nil
}

// ValidateMigrationFile vérifie si un fichier de migration est valide
func (d *MigrationDiscovery) ValidateMigrationFile(filename string) error {
	// Vérifier le format du nom de fichier
	if _, _, err := d.parseMigrationFileName(filename); err != nil {
		return err
	}

	// Vérifier que le fichier existe
	filePath := filepath.Join(d.MigrationsDir, filename)
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("fichier de migration non trouvé: %v", err)
	}

	return nil
}
