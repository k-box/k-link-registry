package mysql

import (
	"strings"

	"github.com/k-box/k-link-registry"
)

// ApplicationRow represents an Application inside the database
type ApplicationRow struct {
	ID          int64  `db:"application_id"`
	OwnerID     int64  `db:"registrant_id"`
	Name        string `db:"name"`
	URL         string `db:"app_domain"`
	Token       string `db:"auth_token"`
	Permissions string `db:"permissions"`
	Active      bool   `db:"status"`
}

func (row *ApplicationRow) fromApplication(app *klinkregistry.Application) {
	if app == nil {
		return
	}

	row.ID = app.ID
	row.OwnerID = app.OwnerID
	row.Name = app.Name
	row.URL = app.URL
	row.Token = app.Token
	row.Permissions = strings.Join(app.Permissions, ",")
	row.Active = app.Active
}

func (row *ApplicationRow) toApplication() *klinkregistry.Application {
	if row == nil {
		return nil
	}

	app := new(klinkregistry.Application)

	app.ID = row.ID
	app.OwnerID = row.OwnerID
	app.Name = row.Name
	app.URL = row.URL
	app.Token = row.Token
	app.Permissions = strings.Split(row.Permissions, ",")
	app.Active = row.Active
	return app
}

// CreateApplication adds a new application inside the database
func (db Database) CreateApplication(app *klinkregistry.Application) error {
	var row ApplicationRow

	row.fromApplication(app)

	res, err := db.db.NamedExec(`INSERT INTO application (
			registrant_id, name, app_domain, auth_token, permissions, status
		) VALUES (
			:registrant_id, :name, :app_domain, :auth_token, :permissions, :status
		)`, &row)
	if err != nil {
		return err
	}

	// Set auto incremented ID
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	row.ID = lastID
	app = row.toApplication() // replace app with the new row information

	return nil
}

// ListApplications returns a list off all applications inside the database
func (db Database) ListApplications() ([]*klinkregistry.Application, error) {
	var rows []*ApplicationRow

	err := db.db.Select(&rows, "SELECT * FROM application ORDER BY application_id ASC")
	if err != nil {
		return nil, err
	}

	var models []*klinkregistry.Application
	for _, row := range rows {
		models = append(models, row.toApplication())
	}

	return models, nil
}

// GetApplicationByID returns a single application by ID
func (db Database) GetApplicationByID(id int64) (*klinkregistry.Application, error) {
	row := new(ApplicationRow)

	err := db.db.Get(row,
		`SELECT * FROM application WHERE application_id=?`,
		id)

	return row.toApplication(), err
}

// GetApplicationByDomain returns a single application by Domain
func (db Database) GetApplicationByDomain(domain string) (*klinkregistry.Application, error) {
	row := new(ApplicationRow)

	err := db.db.Get(row,
		`SELECT * FROM application WHERE app_domain=?`,
		domain)

	return row.toApplication(), err
}

// ReplaceApplication replaces the application inside the dabase, based on
// the ID attribute
func (db Database) ReplaceApplication(app *klinkregistry.Application) error {
	row := new(ApplicationRow)

	row.fromApplication(app)

	_, err := db.db.NamedExec(`UPDATE application SET
		registrant_id = :registrant_id,
		name = :name,
		app_domain = :app_domain,
		auth_token = :auth_token,
		permissions = :permissions,
		status = :status
		WHERE application_id = :application_id`, &row)
	if err != nil {
		return err
	}

	return err
}

// DeleteApplication removes a application entry from the database
func (db Database) DeleteApplication(id int64) error {
	_, err := db.db.Exec("DELETE FROM application WHERE application_id=?", id)
	return err
}
