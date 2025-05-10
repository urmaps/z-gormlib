package gormlib

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MigrationDiscovery gère la découverte automatique des migrations
type MigrationDiscovery struct {
	MigrationsDir string
}

// NewMigrationDiscovery crée un nouveau découvreur de migrations
func NewMigrationDiscovery(migrationsDir string) *MigrationDiscovery {
	return &MigrationDiscovery{
		MigrationsDir: migrationsDir,
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

	var migrations []Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			// Charger le fichier de migration
			migration, err := d.loadMigration(file.Name())
			if err != nil {
				return nil, fmt.Errorf("erreur lors du chargement de la migration %s: %v", file.Name(), err)
			}
			if migration != nil {
				migrations = append(migrations, migration)
			}
		}
	}

	// Trier les migrations par nom (qui contient le timestamp)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name() < migrations[j].Name()
	})

	return migrations, nil
}

// loadMigration charge une migration à partir d'un fichier
func (d *MigrationDiscovery) loadMigration(filename string) (Migration, error) {
	// Construire le chemin complet du fichier
	filePath := filepath.Join(d.MigrationsDir, filename)

	// Lire le contenu du fichier
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du fichier: %v", err)
	}

	// Extraire le nom de la structure de migration
	// Le format attendu est: type MigrationXXXXXXXXXXXX struct{}
	structName := extractStructName(string(content))
	if structName == "" {
		return nil, fmt.Errorf("structure de migration non trouvée dans %s", filename)
	}

	// Créer une instance de la migration
	// Note: Cette partie nécessite que les migrations suivent une convention de nommage
	// et implémentent l'interface Migration
	migration := createMigrationInstance(structName)
	if migration == nil {
		return nil, fmt.Errorf("impossible de créer une instance de la migration %s", structName)
	}

	return migration, nil
}

// extractStructName extrait le nom de la structure de migration du contenu du fichier
func extractStructName(content string) string {
	// Rechercher le pattern "type MigrationXXXXXXXXXXXX struct"
	// où XXXXXXXXXXXX est un timestamp
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "type Migration") && strings.Contains(line, "struct") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "Migration") {
					return part
				}
			}
		}
	}
	return ""
}

// createMigrationInstance crée une instance de migration à partir de son nom
func createMigrationInstance(structName string) Migration {
	// Cette fonction doit être adaptée en fonction de la façon dont vous voulez
	// instancier vos migrations. Une approche courante est d'utiliser un registre
	// de constructeurs de migrations.
	
	// Pour l'instant, nous retournons nil car cette fonction doit être
	// implémentée en fonction de votre architecture spécifique
	return nil
} 