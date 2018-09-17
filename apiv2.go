package klinkregistry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error is a response sent by the API in case of error
type Error struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Context string `json:"context,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

// Errors that the API may emit
var (
	API2ErrDatabase           = Error{500, "Database Error", ""}
	API2ErrEmail              = Error{500, "Error sending Email", ""}
	API2ErrInvalidResponse    = Error{500, "Invalid response", ""}
	API2ErrTokenGeneration    = Error{500, "Token generation error", ""}
	API2ErrUUIDGeneration     = Error{500, "UUID generation failed", ""}
	API2ErrInvalidJSON        = Error{400, "Invalid JSON request", ""}
	API2ErrPasswordRequired   = Error{400, "Password is required", ""}
	API2ErrUnauthorized       = Error{401, "Unauthorized", ""}
	API2ErrInvalidCredentials = Error{403, "Invalid Credentials", ""}
	API2ErrAccountDisabled    = Error{403, "Account disabled", ""}
	API2ErrInvalidURL         = Error{403, "URL could not be understood", ""}
	API2ErrNotFound           = Error{404, "Resource not found", ""}
	API2ErrUserExists         = Error{400, "User already exists", ""}
	API2ErrTokenExpired       = Error{404, "Token has expired", ""}
)

// RegistrationRequest contains all information to start the registtation
// process
type RegistrationRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// PasswordResetRequest contains all information to initiate a password reset
type PasswordResetRequest struct {
	Email string `json:"email"` // current email of the user
}

// EmailChangeRequest contains all information required to initiate a change
// of the primary user email address.
type EmailChangeRequest struct {
	NewEmail string `json:"email"`
}

// API2EmptyResponse can be used for jsonResponse, if the
// API call does not need to return data
type API2EmptyResponse struct{}

// SessionResponse contains information about a user session
type SessionResponse struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

// PermissionModel is the JSON representation of a Permission
type PermissionModel struct {
	Name string `json:"name"`
}

// EmailVerificationModel contains additional information about
// the email verification, for example if a initial password needs to be set
type EmailVerificationModel struct {
	RequirePassword bool   `json:"require_password"`
	DisplayName     string `json:"display_name"`
}

// jsonResponse is a helper function for sending an API Response to the
// client. If response is an `Error` type response, an  error response code
// is sent (which is defined in the `Error.Status`)
func jsonResponse(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var status = 200

	// If the returned object is an error, set status accordingly
	if apiErr, ok := obj.(Error); ok {
		status = apiErr.Status
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		apiErr := API2ErrInvalidResponse
		apiErr.Context = err.Error()
		jsonResponse(w, apiErr)
		return
	}

	w.WriteHeader(status)
	w.Write(bytes)
}

// handleListPermissions provides an endpoint that returns a list of all
// permissions inside the database
func (s *Server) handleListPermissions() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var responses []PermissionModel

		permissions, err := s.store.ListPermissions()
		if s.store.IsNotFound(err) {
			jsonResponse(w, responses)
			return
		} else if err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		for _, permission := range permissions {
			responses = append(responses, PermissionModel(*permission))
		}

		jsonResponse(w, responses)
	}
}

// handlePostRegistration provides an endpoint that allows creation of new
// registrations.
// This handler will create an Registrant and an emailVerification.
func (s *Server) handlePostRegistration() http.HandlerFunc {
	type Response API2EmptyResponse

	return func(w http.ResponseWriter, req *http.Request) {
		var request RegistrationRequest
		var response Response

		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			jsonResponse(w, API2ErrInvalidJSON)
			return
		}

		_, err := s.store.GetRegistrantByEmail(request.Email)
		if err == nil {
			jsonResponse(w, API2ErrUserExists)
			return
		}

		var registrant = &Registrant{
			Name:     request.Name,
			Active:   false,      // Not active until activated by admin
			Password: []byte(""), // login not possible w/ unset pass
			Role:     RoleUser,   // lowest privilege for now
			Email:    request.Email,
		}
		if err := s.store.CreateRegistrant(registrant); err != nil {
			fmt.Println(registrant)
			jsonResponse(w, API2ErrDatabase)
			return
		}

		if err := s.CreateVerification(registrant, registrant.Email); err != nil {
			jsonResponse(w, err)
		}

		// success
		jsonResponse(w, response)
		return
	}
}

func (s *Server) handlePostPasswordReset() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var request PasswordResetRequest
		var response API2EmptyResponse

		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			jsonResponse(w, API2ErrInvalidJSON)
			return
		}

		registrant, err := s.store.GetRegistrantByEmail(request.Email)
		if s.store.IsNotFound(err) {
			// Do not return an API2ErrNotFound, since we dont want to
			// disclose the info that the user does not exist.
			// fake success:
			jsonResponse(w, response)
			return
		} else if err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		// Empty email field means we want to reset PW only.
		if err := s.CreateVerification(registrant, ""); err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		// success
		jsonResponse(w, response)
		return
	}
}
