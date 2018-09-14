package klinkregistry

import (
	"time"

	"github.com/volatiletech/authboss"
	"golang.org/x/crypto/bcrypt"
)

// Registrant contains information about a Registrant. The Email is set after
// registration, but may or may not be confirmed yet. The user is inactive by
// default, until activated by an administrator or owner.
type Registrant struct {
	ID              int64  `db:"registrant_id"`
	Email           string `db:"email"`
	Password        []byte `db:"password"`
	Name            string `db:"name"`
	Role            string `db:"role"`
	Active          bool   `db:"status"`
	LastLogin       int64  `db:"last_login"`
	Confirmed       bool   `db:"confirmed"`
	ConfirmSelector string `db:"confirm_selector"`
	ConfirmVerifier string `db:"confirm_verifier"`
	RecoverSelector string `db:"recover_selector"`
	RecoverVerifier string `db:"recover_verifier"`
	RecoverExpiry   int64  `db:"recover_expiry"`
}

// SetPass sets the value of the password hash to match the provided password.
// The Registrant needs to be saved afterwards to persist the changes.
func (r *Registrant) SetPass(passwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// no errors
	r.Password = hash
	return nil
}

// CheckPass returns true if the provided password matches the stored password
// hash. Always returns false if the password hash is unset.
func (r *Registrant) CheckPass(password string) error {
	err := bcrypt.CompareHashAndPassword(
		r.Password,
		[]byte(password))

	return err
}

// ensure that we've got the right interfaces implemented.
var (
	assertRegistrant = &Registrant{}

	_ authboss.User            = assertRegistrant
	_ authboss.AuthableUser    = assertRegistrant
	_ authboss.ConfirmableUser = assertRegistrant
	_ authboss.RecoverableUser = assertRegistrant
)

// GetPID returns the primary User ID
func (r Registrant) GetPID() string {
	return r.Email
}

// PutPID sets the primary User ID
func (r *Registrant) PutPID(email string) {
	r.Email = email
}

// GetPassword returns the hashed Password of the user
func (r Registrant) GetPassword() string {
	return string(r.Password)
}

// PutPassword stores a hashed password for the user
func (r *Registrant) PutPassword(password string) {
	r.Password = []byte(password)
}

// GetEmail from user
func (r Registrant) GetEmail() string { return r.Email }

// PutEmail into user
func (r *Registrant) PutEmail(email string) { r.Email = email }

// GetConfirmed from user
func (r Registrant) GetConfirmed() bool { return r.Confirmed }

// PutConfirmed into user
func (r *Registrant) PutConfirmed(confirmed bool) { r.Confirmed = confirmed }

// GetConfirmSelector from user
func (r Registrant) GetConfirmSelector() string { return r.ConfirmSelector }

// PutConfirmSelector into user
func (r *Registrant) PutConfirmSelector(confirmSelector string) {
	r.ConfirmSelector = confirmSelector
}

// GetConfirmVerifier from user
func (r Registrant) GetConfirmVerifier() string { return r.ConfirmVerifier }

// PutConfirmVerifier into user
func (r *Registrant) PutConfirmVerifier(confirmVerifier string) {
	r.ConfirmVerifier = confirmVerifier
}

// GetRecoverSelector from user
func (r Registrant) GetRecoverSelector() string { return r.RecoverSelector }

// PutRecoverSelector into user
func (r *Registrant) PutRecoverSelector(selector string) {
	r.RecoverSelector = selector
}

// GetRecoverVerifier from user
func (r Registrant) GetRecoverVerifier() string { return r.RecoverVerifier }

// PutRecoverVerifier into user
func (r *Registrant) PutRecoverVerifier(confirmVerifier string) {
	r.ConfirmVerifier = confirmVerifier
}

// GetRecoverExpiry from user
func (r Registrant) GetRecoverExpiry() time.Time {
	return time.Unix(r.RecoverExpiry, 0)
}

// PutRecoverExpiry into user
func (r *Registrant) PutRecoverExpiry(expiry time.Time) {
	r.RecoverExpiry = expiry.UTC().Unix()
}
