package mysql

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	// MySQL database driver
	_ "github.com/go-sql-driver/mysql"
)

// Database is a MySQL Database
type Database struct {
	db *sqlx.DB
}

// NewDatabase returns a new MySQL database
func NewDatabase(dsn string) (*Database, error) {
	db, err := sqlx.Open("mysql", dsn)

	// Solve the problem of silently dying idle connections
	db.SetConnMaxLifetime(time.Second)

	return &Database{db: db}, err
}

// EmailVerificationRow represents an email verification inside the
// database
type EmailVerificationRow struct {
	Email        string
	RegistrantID int64
	Token        string
	Timestamp    int64
}

// PasswordChangeVerificationRow represents a password change verification
// inside the database
type PasswordChangeVerificationRow struct {
	RegistrantID int64
	Token        string
	Timestamp    int64
}

// PermissionRow represents a Permission inside the database
type PermissionRow struct {
	Name string
}

// IsNotFound returns true, if the error is simply due to no entries being
// found.
func (db Database) IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}
