package klinkregistry

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"github.com/rs/cors"
)

// initRoutes populates the server.routes with a http Handler that serves
// all endpoints, based on the request path.
func (s *Server) initRoutes() {

	apiV1Router := func(r chi.Router) {
		r.Post("/application.authenticate", s.handleAuthenticate())
	}

	apiV2Router := func(r chi.Router) {
		r.Use(s.sessions.Middleware)

		// Authentication Endpoints
		r.Route("/auth", func(r chi.Router) {
			r.Post("/session", s.handleCreateSession())
			r.Get("/session", s.handleGetSession())

			r.Post("/registrations", s.handlePostRegistration())

			r.Get("/email-verification/{token}", s.handleGetVerifyEmail())
			r.Post("/email-verification/{token}", s.handlePostVerifyEmail())

			r.Post("/change-password/{token}", s.handlePostSetPassword())
		})

		// Registrant endpoints
		r.Route("/registrants", func(r chi.Router) {
			r.Use(s.sessions.RequireAuthorized)

			r.Post("/", s.handleCreateRegistrant())
			r.Get("/", s.handleListRegistrants())
			r.Get("/{id}", s.handleGetRegistrant())
			r.Put("/{id}", s.handleUpdateRegistrant())
			r.Delete("/{id}", s.handleDeleteRegistrant())
		})

		// Application endpoints
		r.Route("/applications", func(r chi.Router) {
			r.Use(s.sessions.RequireAuthorized)

			r.Post("/", s.handleCreateApplication())
			r.Get("/", s.handleListApplications())
			r.Get("/{id}", s.handleGetApplication())
			r.Put("/{id}", s.handleUpdateApplication())
			r.Delete("/{id}", s.handleDeleteApplication())
		})

		// K-Links endpoints
		r.Route("/klinks", func(r chi.Router) {
			r.Use(s.sessions.RequireAuthorized)

			r.Post("/", s.handleCreateKlink())
			r.Get("/", s.handleListKlinks())
			r.Get("/{id}", s.handleGetKlink())
			r.Put("/{id}", s.handleUpdateKlink())
			r.Delete("/{id}", s.handleDeleteKlink())
		})

		r.Route("/permissions", func(r chi.Router) {
			r.Get("/", s.handleListPermissions())
		})
	}

	// apiRouter contains all routes for the API endpoints
	apiRouter := func(r chi.Router) {
		// Log all API requests
		r.Use(middleware.Logger)
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
		})
		r.Use(cors.Handler)

		// Define an own not found handler for the API, otherwise it will also
		// serve the UI entrypoint.
		r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "Not found.", http.StatusNotFound)
			return
		})

		// APIs are served with the corresponding version
		r.Route("/1.0", apiV1Router)
		r.Route("/2.0", apiV2Router)
	}

	// baseRouter embeds uses the apiRouter to serve API endpoints, otherwise
	// the static files for the frontend are returned.
	baseRouter := func(r chi.Router) {
		// if no route mathes, serve the UI entrypoint
		r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("x-frame-options", "SAMEORIGIN")
			w.Header().Set("x-content-type-options", "nosniff")
			w.Header().Set("x-xss-protection", "1; mode=block")
			renderFile(w, s.assets, s.config.HTTPBasePath, s.config.NetworkName, "/static/index.html")
			return
		})

		// The {{ .BasePath }} in the manifest file needs to be replaced, so that the static assets
		// can be resolved
		r.HandleFunc("/static/static/js/manifest.{blob}.js", func(w http.ResponseWriter, req *http.Request) {
			blob := chi.URLParam(req, "blob")
			if err := renderFile(w, s.assets, s.config.HTTPBasePath, s.config.NetworkName, "/static/static/js/manifest."+blob+".js"); err != nil {
				fmt.Println(err)
			}
		})

		r.Route("/api", apiRouter)

		r.HandleFunc("/static/*", staticHandler(s.assets, s.config.HTTPBasePath))
	}

	mux := chi.NewMux()

	// ensure router base path contains at least one slash, and is absolute
	//with no trailing slashes
	routerBasePath := path.Join("/", s.config.HTTPBasePath)
	mux.Route(routerBasePath, baseRouter)

	mux.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w,
			`This application does not serve this path, please check the "base path" setting.`,
			http.StatusNotFound,
		)
		return
	})

	s.router = mux
}

// handler functions don’t actually handle the requests, they return
// a function that does. This gives us a closure environment in which
// our handler can operate:
func (s *Server) handleSomething() http.HandlerFunc {
	// Initialization code, will be run once on instanciation, not
	// once per request.
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// staticHandler handles the static assets path.
func staticHandler(fs http.FileSystem, prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if prefix != "/" {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		}
		http.FileServer(fs).ServeHTTP(w, r)
	}
}

// renderFile renders a file using a template with some variables.
func renderFile(w http.ResponseWriter, fs http.FileSystem, baseURL, networkName, path string) error {
	file, err := fs.Open(path)
	if err != nil {
		// could not open file
		return errors.Wrap(err, "Error opening file")
	}

	// we dont usually need to close the file, since we work with virtual
	// assets. For development we will get real files though, so let's try.
	defer file.Close()

	tpl := template.Must(templateFromFile(file))

	var contentType string
	switch filepath.Ext(path) {
	case ".html":
		contentType = "text/html"
	case ".js":
		contentType = "application/javascript"
	case ".json":
		contentType = "application/json"
	case ".yaml", ".yml":
		contentType = "application/yaml"
	default:
		contentType = "text"
	}

	w.Header().Set("Content-Type", contentType+"; charset=utf-8")

	data := map[string]interface{}{
		"BaseURL":     baseURL,
		"NetworkName": networkName,
	}

	err = tpl.Execute(w, data)

	if err != nil {
		return errors.Wrap(err, "Error rendering file")
	}

	return nil
}

// templateFromFile parses the contents of File and returns a template.
func templateFromFile(file http.File) (*template.Template, error) {
	if file == nil {
		return nil, errors.New("Nil file passed in templateFromFile")
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading template file")
	}

	s := string(b)
	return template.New("file").Parse(s)
}
