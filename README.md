# z-gormlib

Une bibliothèque Go pour gérer les migrations de base de données avec GORM, offrant une interface simple et robuste pour la gestion des schémas de base de données.

## Fonctionnalités

- 🚀 Génération automatique de fichiers de migration
- 🔄 Support des migrations et rollbacks
- ⚡ Exécution par lots (batching) des migrations
- 🔒 Transactions pour garantir l'intégrité des données
- ⏱️ Timeouts configurables
- 🔁 Mécanisme de retry automatique
- 📝 Découverte automatique des migrations
- 🔍 Registre de migrations thread-safe
- 🛡️ Validation des noms de migrations
- 🎯 Configuration flexible

## Installation

### Installation de la CLI

```bash
# Installation globale
go install github.com/urmaps/z-gormlib/cmd/gormlib@latest

# Vérifier l'installation
gormlib --help
```

### Installation en tant que dépendance

```bash
go get github.com/urmaps/z-gormlib
```

## Utilisation de la CLI

La bibliothèque fournit une interface en ligne de commande pour gérer les migrations. Voici les commandes disponibles :

```bash
# Afficher l'aide
gormlib --help

# Créer une nouvelle migration
gormlib create create_users_table
# Résultat : Création du fichier migrations/create_users_table.go

# Exécuter les migrations en attente
gormlib migrate
# Résultat : Exécution de toutes les migrations non appliquées

# Annuler la dernière migration
gormlib rollback
# Résultat : Annulation de la dernière migration appliquée
```

### Structure des fichiers générés

Lorsque vous créez une migration avec la CLI, un fichier Go est généré avec la structure suivante :

```go
// migrations/create_users_table.go
package migrations

import "gorm.io/gorm"

type CreateUsersTable struct{}

func (m *CreateUsersTable) Up(db *gorm.DB) error {
    // TODO: Implémenter la migration
    return nil
}

func (m *CreateUsersTable) Down(db *gorm.DB) error {
    // TODO: Implémenter le rollback
    return nil
}

func (m *CreateUsersTable) Name() string {
    return "create_users_table"
}
```

### Configuration de la base de données

La CLI utilise les variables d'environnement suivantes pour la configuration de la base de données :

```bash
# .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=mydb
```

Vous pouvez créer un fichier `.env` à la racine de votre projet pour définir ces variables.

## Utilisation en tant que bibliothèque

### Configuration de la base de données

```go
import "github.com/urmaps/z-gormlib"

// Créer une nouvelle configuration
config := gormlib.NewConfig()

// Configurer la connexion (optionnel)
config.Host = "localhost"
config.Port = 5432
config.User = "postgres"
config.Password = "password"
config.DBName = "mydb"
```

### Configuration

La bibliothèque utilise une configuration par défaut qui peut être personnalisée :

```go
config := &gormlib.MigrationConfig{
    BatchSize:     10,                // Nombre de migrations par lot
    Timeout:       5 * time.Minute,   // Timeout pour les opérations
    RetryAttempts: 3,                 // Nombre de tentatives en cas d'échec
    TableName:     "migrations",      // Nom de la table des migrations
    AutoCreateDir: true,              // Création automatique du dossier migrations
}
```

### Création d'une Migration

```go
// Créer un générateur de migrations
generator := gormlib.NewMigrationGenerator("migrations", gormlib.DefaultConfig())

// Générer une nouvelle migration
err := generator.GenerateMigration("create_users_table")
if err != nil {
    log.Fatal(err)
}
```

### Exécution des Migrations

```go
// Créer un migrator
migrator := gormlib.NewMigrator(db, gormlib.DefaultConfig())

// Créer un registre de migrations
registry := gormlib.NewMigrationRegistry()

// Enregistrer les migrations
registry.Register(&CreateUsersTable{})

// Exécuter les migrations
err := migrator.RunMigrations(registry.GetMigrations()...)
if err != nil {
    log.Fatal(err)
}
```

### Rollback d'une Migration

```go
// Récupérer la dernière migration appliquée
appliedMigrations, err := migrator.GetAppliedMigrations()
if err != nil {
    log.Fatal(err)
}

if len(appliedMigrations) > 0 {
    lastMigration := appliedMigrations[len(appliedMigrations)-1]
    migration, err := registry.GetMigration(lastMigration.Name)
    if err != nil {
        log.Fatal(err)
    }

    // Exécuter le rollback
    err = migrator.RollbackMigration(migration)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Structure d'une Migration

```go
type CreateUsersTable struct{}

func (m *CreateUsersTable) Up(db *gorm.DB) error {
    return db.Migrator().CreateTable(&User{})
}

func (m *CreateUsersTable) Down(db *gorm.DB) error {
    return db.Migrator().DropTable(&User{})
}

func (m *CreateUsersTable) Name() string {
    return "20240321000000_create_users_table"
}
```

## Gestion des Erreurs

La bibliothèque fournit des types d'erreurs spécifiques :

```go
// Erreurs communes
ErrMigrationNotFound     // Migration non trouvée
ErrMigrationAlreadyExists // Migration déjà existante
ErrInvalidMigrationName  // Nom de migration invalide
ErrMigrationFailed      // Échec de la migration
ErrRollbackFailed       // Échec du rollback
```

## Bonnes Pratiques

1. **Nommage des Migrations**
   - Utilisez des noms descriptifs
   - Évitez les espaces et caractères spéciaux
   - Utilisez le format `action_table_name`

2. **Transactions**
   - Les migrations sont exécutées dans des transactions
   - Les rollbacks sont également transactionnels
   - En cas d'échec, la transaction est annulée

3. **Sécurité**
   - Validez toujours les noms de migrations
   - Utilisez des timeouts appropriés
   - Configurez correctement les permissions des fichiers

4. **Performance**
   - Utilisez le batching pour les grandes migrations
   - Configurez le nombre de connexions selon vos besoins
   - Surveillez les timeouts et les retries

## Exemples Complets

### Migration Simple

```go
// Créer une table
type CreateUsersTable struct{}

func (m *CreateUsersTable) Up(db *gorm.DB) error {
    return db.Migrator().CreateTable(&User{})
}

func (m *CreateUsersTable) Down(db *gorm.DB) error {
    return db.Migrator().DropTable(&User{})
}

func (m *CreateUsersTable) Name() string {
    return "20240321000000_create_users_table"
}
```

### Migration Complexe

```go
// Ajouter des colonnes et des index
type AddUserFields struct{}

func (m *AddUserFields) Up(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Ajouter des colonnes
        if err := tx.Migrator().AddColumn(&User{}, "email"); err != nil {
            return err
        }
        if err := tx.Migrator().AddColumn(&User{}, "phone"); err != nil {
            return err
        }
        
        // Créer des index
        return tx.Migrator().CreateIndex(&User{}, "idx_email")
    })
}

func (m *AddUserFields) Down(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Supprimer les index
        if err := tx.Migrator().DropIndex(&User{}, "idx_email"); err != nil {
            return err
        }
        
        // Supprimer les colonnes
        if err := tx.Migrator().DropColumn(&User{}, "phone"); err != nil {
            return err
        }
        return tx.Migrator().DropColumn(&User{}, "email")
    })
}

func (m *AddUserFields) Name() string {
    return "20240321000001_add_user_fields"
}
```

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à :
1. Fork le projet
2. Créer une branche pour votre fonctionnalité
3. Commiter vos changements
4. Pousser vers la branche
5. Ouvrir une Pull Request

## Licence

MIT

### Découverte Automatique des Migrations

GormLib inclut un système de découverte automatique des migrations qui scanne un répertoire pour trouver et charger les migrations disponibles. Ce système élimine le besoin de registre manuel des migrations.

```go
import (
    "github.com/urmaps/z-gormlib"
    "path/filepath"
)

// Créer un découvreur de migrations
discovery := gormlib.NewMigrationDiscovery("path/to/migrations")

// Découvrir toutes les migrations
migrations, err := discovery.DiscoverMigrations()
if err != nil {
    log.Fatalf("Erreur lors de la découverte des migrations: %v", err)
}

// Exécuter les migrations découvertes
migrator := gormlib.NewMigrator(db, config)
if err := migrator.RunMigrations(migrations...); err != nil {
    log.Fatalf("Erreur lors de l'exécution des migrations: %v", err)
}
```

Le système de découverte :
- Scanne récursivement le répertoire spécifié
- Charge automatiquement les fichiers de migration
- Trie les migrations par ordre chronologique
- Gère les erreurs de chargement
- Supporte les migrations Go et SQL

### Création de Migrations

```go
generator := gormlib.NewMigrationGenerator("migrations", config)
err := generator.GenerateMigration("create_users_table")
```

### Exécution des Migrations

```go
migrator := gormlib.NewMigrator(db, config)

// Exécuter toutes les migrations
err := migrator.RunMigrations(migrations...)

// Exécuter une migration spécifique
err := migrator.RunMigration(migration)

// Annuler la dernière migration
err := migrator.RollbackMigration(migration)
```

### Interface en Ligne de Commande

```bash
# Créer une nouvelle migration
gormlib create create_users_table

# Exécuter les migrations
gormlib migrate

# Annuler la dernière migration
gormlib rollback
```

## Structure des Migrations

### Migration Go

```go
package migrations

import (
    "gorm.io/gorm"
)

type CreateUsersTable struct{}

func (m *CreateUsersTable) Up(db *gorm.DB) error {
    return db.AutoMigrate(&User{})
}

func (m *CreateUsersTable) Down(db *gorm.DB) error {
    return db.Migrator().DropTable(&User{})
}

func (m *CreateUsersTable) Name() string {
    return "create_users_table"
}
```

### Migration SQL

```sql
-- 20240321000000_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 20240321000000_create_users_table.down.sql
DROP TABLE users;
```

## Gestion des Erreurs

La bibliothèque fournit des erreurs spécifiques pour différentes situations :

```go
if err != nil {
    switch {
    case errors.Is(err, gormlib.ErrMigrationNotFound):
        // Migration non trouvée
    case errors.Is(err, gormlib.ErrMigrationAlreadyApplied):
        // Migration déjà appliquée
    case errors.Is(err, gormlib.ErrMigrationFailed):
        // Échec de la migration
    case errors.Is(err, gormlib.ErrRollbackFailed):
        // Échec du rollback
    }
}
```

## Meilleures Pratiques

1. **Nommage des Migrations**
   - Utilisez des noms descriptifs
   - Incluez un timestamp dans le nom du fichier
   - Suivez une convention de nommage cohérente

2. **Structure des Migrations**
   - Gardez les migrations atomiques
   - Incluez toujours une méthode `Down`
   - Testez les rollbacks

3. **Organisation**
   - Placez les migrations dans un dossier dédié
   - Utilisez le système de découverte automatique
   - Maintenez un historique des migrations

4. **Sécurité**
   - Ne stockez jamais de données sensibles dans les migrations
   - Utilisez des transactions pour les migrations critiques
   - Sauvegardez la base de données avant les migrations majeures

## Exemples

### Migration Simple

```go
type AddUserRole struct{}

func (m *AddUserRole) Up(db *gorm.DB) error {
    return db.Exec("ALTER TABLE users ADD COLUMN role VARCHAR(50)").Error
}

func (m *AddUserRole) Down(db *gorm.DB) error {
    return db.Exec("ALTER TABLE users DROP COLUMN role").Error
}

func (m *AddUserRole) Name() string {
    return "add_user_role"
}
```

### Migration avec Transaction

```go
type CreateUserTable struct{}

func (m *CreateUserTable) Up(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.AutoMigrate(&User{}); err != nil {
            return err
        }
        return tx.Exec("CREATE INDEX idx_users_email ON users(email)").Error
    })
}

func (m *CreateUserTable) Down(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Exec("DROP INDEX IF EXISTS idx_users_email").Error; err != nil {
            return err
        }
        return tx.Migrator().DropTable(&User{})
    })
}

func (m *CreateUserTable) Name() string {
    return "create_user_table"
}
```

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou une pull request.

## Licence

MIT