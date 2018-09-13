package mysql

import (
	"github.com/k-box/k-link-registry"
)

// CreateRegistrant adds a new Registrant inside the database
func (db Database) CreateRegistrant(r *klinkregistry.Registrant) error {
	res, err := db.db.NamedExec(`INSERT INTO registrant (
			email, password, name, role, status, last_login
		) VALUES (
			:email, :password, :name, :role, :status, :last_login
		)`, &r)

	if err != nil {
		return err
	}

	// Set auto incremented ID
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	r.ID = lastID
	return nil
}

// ListRegistrants returns a list off all registrants inside the database
func (db Database) ListRegistrants() ([]*klinkregistry.Registrant, error) {
	var models []*klinkregistry.Registrant

	err := db.db.Select(&models, "SELECT registrant_id, email, password, name, role, status, last_login FROM registrant ORDER BY registrant_id ASC")
	if err != nil {
		return nil, err
	}

	return models, nil
}

// GetRegistrantByID returns a single registrant by ID
func (db Database) GetRegistrantByID(id int64) (*klinkregistry.Registrant, error) {
	var registrant klinkregistry.Registrant

	err := db.db.Get(&registrant,
		`SELECT registrant_id, email, password, name, role, status, last_login FROM registrant WHERE registrant_id=?`,
		id)

	return &registrant, err
}

// GetRegistrantByEmail returns a single registrant by Email
func (db Database) GetRegistrantByEmail(email string) (*klinkregistry.Registrant, error) {
	var registrant klinkregistry.Registrant

	err := db.db.Get(&registrant,
		`SELECT registrant_id, email, password, name, role, status, last_login FROM registrant WHERE email=?`,
		email)

	return &registrant, err
}

// ReplaceRegistrant replaces the Registrant inside the dabase, based on
// the ID attribute
func (db Database) ReplaceRegistrant(r *klinkregistry.Registrant) error {
	_, err := db.db.NamedExec(`UPDATE registrant SET
		email = :email,
		password = :password,
		name = :name,
		role = :role,
		status = :status,
		last_login = :last_login WHERE registrant_id = :registrant_id`, &r)

	return err
}

// DeleteRegistrant removes a registrant entry from the database
func (db Database) DeleteRegistrant(id int64) error {
	_, err := db.db.Exec("DELETE FROM registrant WHERE registrant_id=?", id)
	return err
}
