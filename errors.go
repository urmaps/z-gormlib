package gormlib

import "fmt"

// MigrationError représente une erreur liée aux migrations
type MigrationError struct {
	Op  string // L'opération qui a échoué
	Err error  // L'erreur sous-jacente
}

func (e *MigrationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

// NewMigrationError crée une nouvelle erreur de migration
func NewMigrationError(op string, err error) error {
	return &MigrationError{Op: op, Err: err}
}

// Common migration errors
var (
	ErrMigrationNotFound     = NewMigrationError("migration not found", nil)
	ErrMigrationAlreadyExists = NewMigrationError("migration already exists", nil)
	ErrInvalidMigrationName  = NewMigrationError("invalid migration name", nil)
	ErrMigrationFailed      = NewMigrationError("migration failed", nil)
	ErrRollbackFailed       = NewMigrationError("rollback failed", nil)
) 