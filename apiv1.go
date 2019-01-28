package klinkregistry

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

// errors
var (
	ErrInvalidPermissions   = errors.New("Invalid Permissions")
	ErrUndefinedApplication = errors.New("Application not defined")
)

// RPCError is returned on failure, contains an error code and an optional
// human readable message
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// RPCResponse wraps all responses for this API version, contains the
// original request id, the result object or an error object
type RPCResponse struct {
	ID     interface{} `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  *RPCError   `json:"error,omitempty"`
}

// KlinkResponse wraps the klink entries in the klinks array
type KlinkResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// writeRPCResponse is a helper function to write the response object to the client
func writeRPCResponse(w http.ResponseWriter, res RPCResponse) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK) // Status is expected to ALWAYS be 200-OK
	return json.NewEncoder(w).Encode(res)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func checkAccess(app *Application, permissions []string) error {
	if app == nil {
		return ErrUndefinedApplication
	}

	for _, permissionName := range permissions {
		if !stringInSlice(permissionName, app.Permissions) {
			return ErrInvalidPermissions
		}
	}
	return nil
}

// MapToKlink maps a list of ids to the corresponding K-Link instance
func (s *Server) MapToKlink(vs []string) []KlinkResponse {
	vsm := make([]KlinkResponse, 0)
	for _, v := range vs {

		if v != "" {
			klink, err := s.store.GetKlinkByIdentifier(v)

			if err == nil {
				model := new(KlinkResponse)
				model.ID = klink.Identifier
				model.Name = klink.Name
				vsm = append(vsm, *model)
			}
		}
	}
	return vsm
}

// handleAuthenticate will serve the application.authenticate endpoint on
// the v1 API.
func (s *Server) handleAuthenticate() http.HandlerFunc {

	// Request contains the request as expected from the client
	type Request struct {
		ID         interface{} `json:"id"`
		Parameters struct {
			AppSecret   string   `json:"app_secret"`
			AppURL      string   `json:"app_url"`
			Permissions []string `json:"permissions"`
		} `json:"params"`
	}

	// AppResponse is the Application returned to the client on success
	type AppResponse struct {
		Name        string          `json:"name"`
		AppURL      string          `json:"app_url"`
		AppID       int64           `json:"app_id"`
		Permissions []string        `json:"permissions"`
		Klinks      []KlinkResponse `json:"klinks"`
		OwnerEmail  string          `json:"email"`
	}

	// Common errors we will encounter
	var (
		APIErrInvalidJSON      = RPCError{-32700, "Invalid JSON object."}
		APIErrInvalidRequest   = RPCError{-32602, "Invalid request."}
		APIErrPermissionDenied = RPCError{-32000, "Permission Denied."}
	)

	validate := validator.New()

	// returned anonymous function contains the actual handler function
	return func(w http.ResponseWriter, req *http.Request) {
		var response = RPCResponse{}
		var request = Request{}
		var decoder = json.NewDecoder(req.Body)

		if err := decoder.Decode(&request); err != nil {
			response.Error = &APIErrInvalidJSON
			log.Println("v1-application validation malformed request payload")
			writeRPCResponse(w, response)
			return
		}

		log.Printf("v1-application validation [%s]", request.Parameters.AppURL)

		// ResponseID should be the same as the request ID
		response.ID = request.ID

		if err := validate.Struct(request); err != nil {
			response.Error = &APIErrInvalidRequest
			writeRPCResponse(w, response)
			return
		}

		app, err := s.store.GetApplicationByDomain(request.Parameters.AppURL)
		if err != nil {
			response.Error = &APIErrPermissionDenied
			writeRPCResponse(w, response)
			return
		}

		// Check if the secret token matches
		if app.Token != request.Parameters.AppSecret {
			response.Error = &APIErrPermissionDenied
			writeRPCResponse(w, response)
			return
		}

		if err := checkAccess(app, request.Parameters.Permissions); err != nil {
			response.Error = &APIErrPermissionDenied
			writeRPCResponse(w, response)
			return
		}

		// fetch the Application owner
		owner, error := s.store.GetRegistrantByID(app.OwnerID)
		if error != nil {
			response.Error = &APIErrPermissionDenied
			writeRPCResponse(w, response)
			return
		}

		response.Result = AppResponse{
			Name:        app.Name,
			AppURL:      app.URL,
			AppID:       app.ID,
			Permissions: app.Permissions,
			Klinks:      s.MapToKlink(app.Klinks),
			OwnerEmail:  owner.Email,
		}

		writeRPCResponse(w, response)
	}
}
