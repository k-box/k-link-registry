package mysql

import (
	"log"

	klinkregistry "github.com/k-box/k-link-registry"
)

// KlinkRow represents an Klink inside the database
type KlinkRow struct {
	ID          int64  `db:"klink_id"`
	Identifier  string `db:"identifier"`
	ManagerID   int64  `db:"manager_id"`
	Name        string `db:"name"`
	Website     string `db:"website"`
	Description string `db:"description"`
	Active      bool   `db:"active"`
}

func (row *KlinkRow) fromKlink(klink *klinkregistry.Klink) *KlinkRow {
	if klink == nil {
		return nil
	}

	row.ID = klink.ID
	row.Identifier = klink.Identifier
	row.ManagerID = klink.ManagerID
	row.Name = klink.Name
	row.Website = klink.Website
	row.Description = klink.Description
	row.Active = klink.Active

	return row
}

func (row *KlinkRow) toKlink() *klinkregistry.Klink {
	if row == nil {
		return nil
	}

	app := new(klinkregistry.Klink)

	app.ID = row.ID
	app.Identifier = row.Identifier
	app.ManagerID = row.ManagerID
	app.Name = row.Name
	app.Website = row.Website
	app.Description = row.Description
	app.Active = row.Active
	return app
}

// CreateKlink adds a new klink inside the database
func (db Database) CreateKlink(app *klinkregistry.Klink) error {
	var row KlinkRow

	row.fromKlink(app)

	res, err := db.db.NamedExec(`INSERT INTO klink (
			identifier, manager_id, name, website, description, active
		) VALUES (
			:identifier, :manager_id, :name, :website, :description,  :active
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
	app = row.toKlink() // replace app with the new row information

	return nil
}

// ListKlinks returns a list off all klinks inside the database
func (db Database) ListKlinks() ([]*klinkregistry.Klink, error) {
	var rows []*KlinkRow

	err := db.db.Select(&rows, "SELECT * FROM klink ORDER BY klink_id ASC")
	if err != nil {
		return nil, err
	}

	var models []*klinkregistry.Klink
	for _, row := range rows {
		models = append(models, row.toKlink())
	}

	return models, nil
}

// GetKlinkByPrimaryKey returns a single klink by ID
func (db Database) GetKlinkByPrimaryKey(id int64) (*klinkregistry.Klink, error) {
	row := new(KlinkRow)

	err := db.db.Get(row,
		`SELECT * FROM klink WHERE klink_id=?`,
		id)

	return row.toKlink(), err
}

// GetKlinkByIdentifier returns a single klink by its public identifier
func (db Database) GetKlinkByIdentifier(id string) (*klinkregistry.Klink, error) {
	row := new(KlinkRow)

	err := db.db.Get(row,
		`SELECT * FROM klink WHERE identifier=?`,
		id)

	return row.toKlink(), err
}

// UpdateKlink update the klink inside the dabase, based on
// the ID attribute
func (db Database) UpdateKlink(app *klinkregistry.Klink) error {
	row := new(KlinkRow)

	row.fromKlink(app)

	_, err := db.db.NamedExec(`UPDATE klink SET
		manager_id = :manager_id,
		name = :name,
		website = :website,
		description = :description,
		active = :active
		WHERE identifier = :identifier`, &row)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

// DeleteKlink removes a klink entry from the database
func (db Database) DeleteKlink(id int64) error {
	_, err := db.db.Exec("DELETE FROM klink WHERE identifier=?", id)
	return err
}
