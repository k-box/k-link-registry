package klinkregistry

import (
	"time"
)

// Application contains information about a registered Application.
type Application struct {
	ID          int64    `db:"application_id"`
	OwnerID     int64    `db:"registrant_id"`
	Name        string   `db:"name"`
	URL         string   `db:"app_domain"`
	Token       string   `db:"auth_token"`
	Permissions []string `db:"permissions"`
	Active      bool     `db:"status"`
}

// EmailVerification represents a emailVerification in the database On
// registration, an first email verification is created. After the email is
// verified for a user with no password, the user will also be asked to set a
// password.
type EmailVerification struct {
	Email        string `db:"email"`
	RegistrantID int64  `db:"registrant_id"`
	Token        string `db:"token"`
	Timestamp    int64  `db:"timestamp"`
}

// IsExpired returns true if the email verification is expired
func (v EmailVerification) IsExpired() bool {
	now := time.Now().UTC().Unix()

	if v.Timestamp+60*60*24 > now {
		return true
	}
	return false
}

// PasswordChangeVerification represents a password change verification in
// the database.
type PasswordChangeVerification struct {
	RegistrantID int64  `db:"registrant_id"`
	Token        string `db:"token"`
	Timestamp    int64  `db:"timestamp"`
}

// A Permission describes an action that an Application may perform
type Permission struct {
	Name string `db:"name"`
}

// A EmailConfirmation represents a token (sent via email) that a registrant may
// use to confirm that he has access to the claimed mailbox
type EmailConfirmation struct {
	ID               int64     `db:"id"`
	Token            string    `db:"token"`
	RegistrantID     int64     `db:"registrant_id"`
	CreatedAt        time.Time `db:"created_at"`
	ValidUntil       time.Time `db:"valid_until"`
	ForceSetPassword bool      `db:"force_set_password"`
	NewAddress       string    `db:"new_address"`
}

// A PasswordReset represents a token (sent via email) that a registrant can use
// to change the current password. This is either due to the user forgetting or
// updating their password.
type PasswordReset struct {
	ID           int64     `db:"id"`
	Token        string    `db:"token"`
	RegistrantID int64     `db:"registrant_id"`
	CreatedAt    time.Time `db:"created_at"`
	ValidUntil   time.Time `db:"valid_until"`
}
