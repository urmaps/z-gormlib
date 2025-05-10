package gormlib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	migrationTemplate = `package migrations

import (
	"gorm.io/gorm"
)

// %s représente la migration %s
type %s struct{}

// Up effectue la migration
func (m *%s) Up(db *gorm.DB) error {
	// TODO: Ajoutez vos modifications ici
	// Exemple:
	// return db.Migrator().CreateTable(&YourModel{})
	// ou pour ajouter une colonne:
	// return db.Migrator().AddColumn(&YourModel{}, "new_column")
	return nil
}

// Down effectue le rollback
func (m *%s) Down(db *gorm.DB) error {
	// TODO: Ajoutez vos rollbacks ici
	// Exemple:
	// return db.Migrator().DropTable(&YourModel{})
	// ou pour supprimer une colonne:
	// return db.Migrator().DropColumn(&YourModel{}, "new_column")
	return nil
}

// Name retourne le nom de la migration
func (m *%s) Name() string {
	return "%s"
}
`
)

// MigrationGenerator gère la génération des fichiers de migration
type MigrationGenerator struct {
	MigrationsDir string
	config        *MigrationConfig
}

// NewMigrationGenerator crée un nouveau générateur de migrations
func NewMigrationGenerator(migrationsDir string, config *MigrationConfig) *MigrationGenerator {
	if config == nil {
		config = DefaultConfig()
	}
	return &MigrationGenerator{
		MigrationsDir: migrationsDir,
		config:        config,
	}
}

// GenerateMigration crée une nouvelle migration à partir d'un nom
func (g *MigrationGenerator) GenerateMigration(name string) error {
	// Valider le nom de la migration
	if err := g.validateMigrationName(name); err != nil {
		return err
	}

	// Créer le timestamp pour le nom du fichier
	timestamp := time.Now().Format("20060102150405")
	migrationName := fmt.Sprintf("%s_%s", timestamp, strings.ToLower(name))
	
	// Créer le nom de la structure Go
	structName := fmt.Sprintf("%s%s", MigrationStructPrefix, strings.Title(name))

	// Créer le dossier migrations s'il n'existe pas
	if g.config.AutoCreateDir {
		if err := os.MkdirAll(g.MigrationsDir, DefaultDirMode); err != nil {
			return NewMigrationError("create migrations directory", err)
		}
	}

	// Vérifier si le fichier existe déjà
	filePath := filepath.Join(g.MigrationsDir, fmt.Sprintf("%s%s", migrationName, MigrationFileSuffix))
	if _, err := os.Stat(filePath); err == nil {
		return ErrMigrationAlreadyExists
	}

	// Créer le fichier de migration
	content := fmt.Sprintf(migrationTemplate, 
		structName, name, structName, structName, structName, structName, migrationName)
	
	if err := os.WriteFile(filePath, []byte(content), DefaultFileMode); err != nil {
		return NewMigrationError("create migration file", err)
	}

	return nil
}

// validateMigrationName vérifie si le nom de la migration est valide
func (g *MigrationGenerator) validateMigrationName(name string) error {
	if name == "" {
		return ErrInvalidMigrationName
	}

	// Vérifier que le nom ne contient que des caractères autorisés
	if !isValidMigrationName(name) {
		return ErrInvalidMigrationName
	}

	return nil
}

// isValidMigrationName vérifie si le nom de la migration est valide
func isValidMigrationName(name string) bool {
	// Le nom ne doit contenir que des lettres, des chiffres et des underscores
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
} 