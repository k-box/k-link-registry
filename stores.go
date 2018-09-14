package klinkregistry

import (
	"github.com/volatiletech/authboss"
)

// RegistrantStorer implements all methods to persist Registrants
type RegistrantStorer interface {
	CreateRegistrant(*Registrant) error
	ListRegistrants() ([]*Registrant, error)
	GetRegistrantByID(id int64) (*Registrant, error)
	GetRegistrantByEmail(email string) (*Registrant, error)
	ReplaceRegistrant(u *Registrant) error
	DeleteRegistrant(id int64) error
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

// PermissionStorer implements all methods to persist Permissions
type PermissionStorer interface {
	ListPermissions() ([]*Permission, error)
	CreatePermission(*Permission) error
}

// A Storer implements all neccessary database methods
type Storer interface {
	authboss.ServerStorer
	RegistrantStorer
	ApplicationStorer
	PermissionStorer
	IsNotFound(error) bool
}
