# z-gormlib

Une bibliothèque Go pour gérer les migrations de base de données avec GORM.

## Table des Matières

- [Fonctionnalités](#fonctionnalités)
- [Installation](#installation)
- [Utilisation Rapide](#utilisation-rapide)
- [Guide d'Utilisation](#guide-dutilisation)
  - [Configuration](#configuration)
  - [Création de Migrations](#création-de-migrations)
  - [Exécution des Migrations](#exécution-des-migrations)
  - [Rollback](#rollback)
- [Interface en Ligne de Commande](#interface-en-ligne-de-commande)
- [Meilleures Pratiques](#meilleures-pratiques)
- [Exemples](#exemples)
- [Contribution](#contribution)
- [Licence](#licence)

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

```bash
go get github.com/urmaps/z-gormlib
```

## Utilisation Rapide

```go
import "github.com/urmaps/z-gormlib"

// Configuration
config := gormlib.DefaultConfig()

// Connexion à la base de données
conn, err := gormlib.NewConnection(config)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Créer un migrator
migrator := gormlib.NewMigrator(conn.DB(), config)

// Découvrir et exécuter les migrations
discovery := gormlib.NewMigrationDiscovery("migrations")
migrations, err := discovery.DiscoverMigrations()
if err != nil {
    log.Fatal(err)
}

if err := migrator.RunMigrations(migrations...); err != nil {
    log.Fatal(err)
}
```

## Guide d'Utilisation

### Configuration

```go
// Configuration par défaut
config := gormlib.DefaultConfig()

// Configuration personnalisée
config := gormlib.NewConfig()
config.SetTableName("custom_migrations")
config.SetLockTimeout(30)
```

### Création de Migrations

```go
// Créer un générateur de migrations
generator := gormlib.NewMigrationGenerator("migrations", config)

// Générer une nouvelle migration
err := generator.GenerateMigration("create_users_table")
```

### Exécution des Migrations

```go
// Exécuter toutes les migrations
err := migrator.RunMigrations(migrations...)

// Exécuter une migration spécifique
err := migrator.RunMigration(migration)
```

### Rollback

```go
// Annuler la dernière migration
err := migrator.RollbackMigration(migration)
```

## Interface en Ligne de Commande

```bash
# Créer une nouvelle migration
gormlib -create-migration create_users_table

# Exécuter les migrations
gormlib -migrate

# Annuler la dernière migration
gormlib -rollback

# Spécifier un dossier de migrations
gormlib -dir custom/migrations -migrate
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

### Migration avec Transaction

```go
type AddUserRole struct{}

func (m *AddUserRole) Up(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Exec("ALTER TABLE users ADD COLUMN role VARCHAR(50)").Error; err != nil {
            return err
        }
        return tx.Exec("CREATE INDEX idx_users_role ON users(role)").Error
    })
}

func (m *AddUserRole) Down(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Exec("DROP INDEX IF EXISTS idx_users_role").Error; err != nil {
            return err
        }
        return tx.Exec("ALTER TABLE users DROP COLUMN role").Error
    })
}

func (m *AddUserRole) Name() string {
    return "add_user_role"
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