package mysql

import (
	"github.com/k-box/k-link-registry"
)

// CreateEmailVerification adds a new EmailVerification inside the database
func (db Database) CreateEmailVerification(r *klinkregistry.EmailVerification) error {
	_, err := db.db.NamedExec(`INSERT INTO email_verification (
			email, registrant_id, token, timestamp
		) VALUES (
			:email, :registrant_id, :token, :timestamp
		)`, &r)

	return err
}

// GetEmailVerificationByRegistrantID returns a single EmailVerification by Registrant ID
func (db Database) GetEmailVerificationByRegistrantID(id int64) (*klinkregistry.EmailVerification, error) {
	var model klinkregistry.EmailVerification

	err := db.db.Get(&model,
		`SELECT email, registrant_id, token, timestamp FROM email_verification WHERE registrant_id=?`,
		id)

	return &model, err
}

// GetEmailVerificationByToken returns a single EmailVerification by Token
func (db Database) GetEmailVerificationByToken(token string) (*klinkregistry.EmailVerification, error) {
	var model klinkregistry.EmailVerification

	err := db.db.Get(&model,
		`SELECT email, registrant_id, token, timestamp FROM email_verification WHERE token=?`,
		token)

	return &model, err
}

// DeleteEmailVerification removes a EmailVerification entry from the database
func (db Database) DeleteEmailVerification(registrantID int64) error {
	_, err := db.db.Exec("DELETE FROM email_verification WHERE registrant_id=?", registrantID)
	return err
}
