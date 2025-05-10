package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urmaps/z-gormlib"
)

func main() {
	// Définir les sous-commandes
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	rollbackCmd := flag.NewFlagSet("rollback", flag.ExitOnError)

	// Vérifier si une sous-commande a été fournie
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Traiter les sous-commandes
	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if createCmd.NArg() == 0 {
			fmt.Println("Erreur: nom de la migration requis")
			createCmd.PrintDefaults()
			os.Exit(1)
		}
		handleCreate(createCmd.Arg(0))

	case "migrate":
		migrateCmd.Parse(os.Args[2:])
		handleMigrate()

	case "rollback":
		rollbackCmd.Parse(os.Args[2:])
		handleRollback()

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  gormlib create <migration_name>  Créer une nouvelle migration")
	fmt.Println("  gormlib migrate                  Exécuter les migrations en attente")
	fmt.Println("  gormlib rollback                 Annuler la dernière migration")
}

func handleCreate(name string) {
	// Définir le dossier des migrations
	migrationsDir := "migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		log.Fatalf("Erreur lors de la création du dossier migrations: %v", err)
	}

	// Créer le générateur de migration
	generator := gormlib.NewMigrationGenerator(migrationsDir)

	// Générer la migration
	if err := generator.GenerateMigration(name); err != nil {
		log.Fatalf("Erreur lors de la création de la migration: %v", err)
	}

	fmt.Printf("Migration '%s' créée avec succès dans le dossier '%s'\n", name, migrationsDir)
}

func handleMigrate() {
	// Configuration de la base de données
	config := gormlib.NewConfig()
	
	// Connexion à la base de données
	conn, err := gormlib.NewConnection(config)
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}
	defer conn.Close()

	// Création du migrator
	migrator := gormlib.NewMigrator(conn.GetDB())

	// Création du registre de migrations
	registry := gormlib.NewMigrationRegistry()

	// Découverte automatique des migrations
	migrationsDir := "migrations"
	if err := filepath.Walk(migrationsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// TODO: Charger dynamiquement les migrations
			return nil
		}
		return nil
	}); err != nil {
		log.Fatalf("Erreur lors de la découverte des migrations: %v", err)
	}

	// Exécution des migrations
	migrations := registry.GetMigrations()
	if err := migrator.RunMigrations(migrations...); err != nil {
		log.Fatalf("Erreur lors de l'exécution des migrations: %v", err)
	}

	fmt.Println("Migrations exécutées avec succès")
}

func handleRollback() {
	// Configuration de la base de données
	config := gormlib.NewConfig()
	
	// Connexion à la base de données
	conn, err := gormlib.NewConnection(config)
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}
	defer conn.Close()

	// Création du migrator
	migrator := gormlib.NewMigrator(conn.GetDB())

	// Récupérer les migrations appliquées
	appliedMigrations, err := migrator.GetAppliedMigrations()
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des migrations appliquées: %v", err)
	}

	if len(appliedMigrations) == 0 {
		fmt.Println("Aucune migration à annuler")
		return
	}

	// Récupérer la dernière migration appliquée
	lastMigration := appliedMigrations[len(appliedMigrations)-1]
	
	// Création du registre de migrations
	registry := gormlib.NewMigrationRegistry()

	// Récupérer la migration correspondante
	migration, err := registry.GetMigration(lastMigration.Name)
	if err != nil {
		log.Fatalf("Migration %s non trouvée: %v", lastMigration.Name, err)
	}

	// Exécuter le rollback
	if err := migrator.RollbackMigration(migration); err != nil {
		log.Fatalf("Erreur lors du rollback de la migration: %v", err)
	}

	fmt.Printf("Migration %s annulée avec succès\n", lastMigration.Name)
} 