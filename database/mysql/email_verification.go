package mysql

import (
	"git.klink.asia/main/klinkregistry"
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

// GetEmailVerificationByEmail returns a single EmailVerification by Email
func (db Database) GetEmailVerificationByEmail(email string) (*klinkregistry.EmailVerification, error) {
	var model klinkregistry.EmailVerification

	err := db.db.Get(&model,
		`SELECT email, registrant_id, token, timestamp FROM email_verification WHERE email=?`,
		email)

	return &model, err
}

// GetEmailVerificationByToken returns a single EmailVerification by Email
func (db Database) GetEmailVerificationByToken(token string) (*klinkregistry.EmailVerification, error) {
	var model klinkregistry.EmailVerification

	err := db.db.Get(&model,
		`SELECT email, registrant_id, token, timestamp FROM email_verification WHERE token=?`,
		token)

	return &model, err
}

// DeleteEmailVerification removes a EmailVerification entry from the database
func (db Database) DeleteEmailVerification(email string) error {
	_, err := db.db.Exec("DELETE FROM email_verification WHERE email=?", email)
	return err
}
