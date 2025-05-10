package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urmaps/z-gormlib"
)

func main() {
	// Charger les variables d'environnement
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Définir les flags
	createMigration := flag.String("create-migration", "", "Create a new migration with the specified name")
	migrate := flag.Bool("migrate", false, "Exécute les migrations en attente")
	rollback := flag.Bool("rollback", false, "Annule la dernière migration")
	migrationsDir := flag.String("dir", "migrations", "Directory containing migrations")
	flag.Parse()

	// Configuration par défaut
	config := gormlib.DefaultConfig()

	// Si on veut créer une migration, on le fait avant de se connecter à la base de données
	if *createMigration != "" {
		generator := gormlib.NewMigrationGenerator(*migrationsDir, config)
		if err := generator.GenerateMigration(*createMigration); err != nil {
			log.Fatalf("Erreur lors de la création de la migration: %v", err)
		}
		return
	}

	// Configuration de la base de données
	dbConfig := gormlib.NewConfig()
	
	// Connexion à la base de données
	conn, err := gormlib.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}
	defer conn.Close()

	// Création du migrator
	migrator := gormlib.NewMigrator(conn.DB(), config)

	// Création du découvreur de migrations
	discovery := gormlib.NewMigrationDiscovery(*migrationsDir)

	if *migrate {
		// Découvrir et exécuter les migrations
		migrations, err := discovery.DiscoverMigrations()
		if err != nil {
			log.Fatalf("Erreur lors de la découverte des migrations: %v", err)
		}
		
		if err := migrator.RunMigrations(migrations...); err != nil {
			log.Fatalf("Erreur lors de l'exécution des migrations: %v", err)
		}
		fmt.Println("Migrations exécutées avec succès")
		return
	}

	if *rollback {
		// Récupérer les migrations appliquées
		appliedMigrations, err := migrator.GetAppliedMigrations()
		if err != nil {
			log.Fatalf("Erreur lors de la récupération des migrations appliquées: %v", err)
		}

		if len(appliedMigrations) == 0 {
			log.Println("Aucune migration à annuler")
			return
		}

		// Récupérer la dernière migration appliquée
		lastMigration := appliedMigrations[len(appliedMigrations)-1]
		
		// Découvrir toutes les migrations
		migrations, err := discovery.DiscoverMigrations()
		if err != nil {
			log.Fatalf("Erreur lors de la découverte des migrations: %v", err)
		}

		// Trouver la migration correspondante
		var targetMigration gormlib.Migration
		for _, m := range migrations {
			if m.Name() == lastMigration.Name {
				targetMigration = m
				break
			}
		}

		if targetMigration == nil {
			log.Fatalf("Migration %s non trouvée", lastMigration.Name)
		}

		// Exécuter le rollback
		if err := migrator.RollbackMigration(targetMigration); err != nil {
			log.Fatalf("Erreur lors du rollback de la migration: %v", err)
		}
		log.Printf("Migration %s annulée avec succès", lastMigration.Name)
		return
	}

	// Si aucun flag n'est spécifié, afficher l'aide
	flag.Usage()
	os.Exit(1)
} 