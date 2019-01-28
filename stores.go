package klinkregistry

// RegistrantStorer implements all methods to persist Registrants
type RegistrantStorer interface {
	CreateRegistrant(*Registrant) error
	ListRegistrants() ([]*Registrant, error)
	GetRegistrantByID(id int64) (*Registrant, error)
	GetRegistrantByEmail(email string) (*Registrant, error)
	ReplaceRegistrant(u *Registrant) error
	DeleteRegistrant(id int64) error
}

// EmailVerificationStorer implements all methods to persist email verifications
type EmailVerificationStorer interface {
	CreateEmailVerification(*EmailVerification) error
	GetEmailVerificationByEmail(string) (*EmailVerification, error)
	GetEmailVerificationByToken(string) (*EmailVerification, error)
	DeleteEmailVerification(email string) error
}

// ApplicationStorer implements all methods to persist Applications
type ApplicationStorer interface {
	CreateApplication(*Application) error
	ListApplications() ([]*Application, error)
	GetApplicationByID(id int64) (*Application, error)
	GetApplicationByDomain(domain string) (*Application, error)
	ReplaceApplication(*Application) error
	DeleteApplication(id int64) error
}

// KlinkStorer implements all methods to persist Klinks
type KlinkStorer interface {
	CreateKlink(*Klink) error
	ListKlinks() ([]*Klink, error)
	GetKlinkByPrimaryKey(id int64) (*Klink, error)
	GetKlinkByIdentifier(identifier string) (*Klink, error)
	UpdateKlink(*Klink) error
	DeleteKlink(id int64) error
}

// PermissionStorer implements all methods to persist Permissions
type PermissionStorer interface {
	ListPermissions() ([]*Permission, error)
	CreatePermission(*Permission) error
}

// A Storer implements all neccessary database methods
type Storer interface {
	RegistrantStorer
	ApplicationStorer
	PermissionStorer
	EmailVerificationStorer
	KlinkStorer
	IsNotFound(error) bool
}
