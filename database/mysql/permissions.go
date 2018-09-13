package mysql

import "github.com/k-box/k-link-registry"

// ListPermissions returns a list off all permissions inside the database
func (db Database) ListPermissions() ([]*klinkregistry.Permission, error) {
	var models []*klinkregistry.Permission

	err := db.db.Select(&models, "SELECT * FROM permission ORDER BY name")
	if err != nil {
		return nil, err
	}

	return models, nil
}

// CreatePermission adds a new Permission inside the database
func (db Database) CreatePermission(p *klinkregistry.Permission) error {
	_, err := db.db.NamedExec(`INSERT INTO permission (
		name
	) VALUES (
		:name
	)`, &p)

	if err != nil {
		return err
	}
	return nil
}
