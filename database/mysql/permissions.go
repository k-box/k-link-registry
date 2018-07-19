package mysql

import "git.klink.asia/main/klinkregistry"

// ListPermissions returns a list off all permissions inside the database
func (db Database) ListPermissions() ([]*klinkregistry.Permission, error) {
	var models []*klinkregistry.Permission

	err := db.db.Select(&models, "SELECT * FROM permission ORDER BY name")
	if err != nil {
		return nil, err
	}

	return models, nil
}
