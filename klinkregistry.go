package klinkregistry

import (
	"crypto/rand"
	"math"
	"net/http"
	"time"

	"git.klink.asia/main/klinkregistry/assets"
	"git.klink.asia/main/klinkregistry/mail"
	"github.com/pkg/errors"
)

// Version will be set automatically on release builds, using a build
// command such as:
//     go build -ldflags "-X main.Version=`git rev-parse --abbrev-ref HEAD`"
var Version = "development (untagged)"

// A Emailer is an interface that can send emails to users
type Emailer interface {
	Email(recepient, subject, html, text string) error
}

// Config contains the configuration for the Application
type Config struct {
	AssetDir string // use embedded assets if empty

	NetworkName string // Name of the managed Network

	HTTPListen       string
	HTTPMaxHeader    int
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPDomain       string
	HTTPBasePath     string
	HTTPSecret       string

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	AdminUsername string
	AdminPassword string
}

// Server is a struct that serves the Web application
type Server struct {
	assets   http.FileSystem
	router   http.Handler
	email    Emailer
	store    Storer
	config   *Config
	sessions SessionsProvider
}

// SetStore is a setter for setting a database inside the application.
// this was originally named `initDatabase`, but caused a cyclic dependency
// FIXME
func (s *Server) SetStore(store Storer) error {
	s.store = store
	return nil
}

func (s *Server) initSMTP() error {
	if s.config.SMTPHost == "" {
		// Init debug mailer if hostname is empty
		s.email = &mail.DebugMailer{}
		return nil
	}

	if s.config.SMTPPort > math.MaxUint16 || s.config.SMTPPort < 1 {
		return errors.New("SMTP Port out of range")
	}

	smtpOptions := &mail.SMTPOptions{
		Host: s.config.SMTPHost,
		Port: s.config.SMTPPort,
		User: s.config.SMTPUser,
		Pass: s.config.SMTPPassword,
		From: s.config.SMTPFrom,
	}
	s.email = &mail.SMTPMailer{Options: smtpOptions}

	return nil
}

// NewServer initializes a new server from the configuration
func NewServer(config *Config) (*Server, error) {
	s := &Server{}
	s.config = config

	// if no assets dir is specified, use the internally packaged assets.
	// otherwise initialize the external assets file.
	if s.config.AssetDir == "" {
		s.assets = assets.Assets
	} else {
		s.assets = http.Dir(s.config.AssetDir)
	}

	s.sessions = &JWTSession{Key: []byte(s.config.HTTPSecret)}
	// use a random key if an empty one was set
	if s.config.HTTPSecret == "" {
		key := make([]byte, 64)
		rand.Read(key)
		s.sessions = &JWTSession{Key: key}
	}

	s.initSMTP()
	s.initRoutes()

	return s, nil
}

// Run starts serving of the HTTP Endpoints
func (s Server) Run() error {
	if s.router == nil {
		panic("Router not initialized")
	}

	readTimeout := 10 * time.Second
	writeTimeout := readTimeout

	server := http.Server{
		Addr:           s.config.HTTPListen,
		Handler:        s.router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}
