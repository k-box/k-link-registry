package mysql

import (
	"database/sql"
	"net/http"

	vfs "git.klink.asia/paul/migrate-vfs"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

// GetMigrator returns a Database Migrator for MySQL. Supported mathods are
// e.g.:
//  * `Up` – Migrate to the latest version
//  * `Down` – empty everything
//  * `(integer)` – Migrate to specific version
func GetMigrator(db *sql.DB, fs http.FileSystem, path string) (*migrate.Migrate, error) {
	source, err := vfs.WithInstance(fs, path)
	if err != nil {
		return nil, err
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, err
	}

	// the strings are only for logging purpose
	migrator, err := migrate.NewWithInstance(
		"vfs-dir", source,
		"mysql", driver,
	)
	if err != nil {
		return nil, err
	}

	return migrator, err
}
