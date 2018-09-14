package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"

	klinkregistry "github.com/k-box/k-link-registry"
	"github.com/volatiletech/authboss"

	"github.com/jmoiron/sqlx"

	// MySQL database driver
	_ "github.com/go-sql-driver/mysql"
)

// ensure that we've got the right interfaces implemented.
var (
	assertDatabase = &Database{}

	_ authboss.ServerStorer = assertDatabase
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

// PermissionRow represents a Permission inside the database
type PermissionRow struct {
	Name string
}

// IsNotFound returns true, if the error is simply due to no entries being
// found.
func (db Database) IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}

// Load user from the database, based on session data
func (db Database) Load(ctx context.Context, key string) (authboss.User, error) {
	return db.GetRegistrantByEmail(key)
}

// Save user to the database, with PID as PK.
func (db Database) Save(ctx context.Context, user authboss.User) error {
	// Check type of user object
	switch ut := user.(type) {
	case *klinkregistry.Registrant:
		if ut.ID != 0 {
			// is not a fresh user,
			return authboss.ErrUserFound
		}
		return db.ReplaceRegistrant(ut)
	default:
		return errors.Errorf("User has invalid type %T", ut)
	}
}
