package klinkregistry

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

// Error is a response sent by the API in case of error
type Error struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Context string `json:"context,omitempty"`
}

// Errors that the API may emit
var (
	API2ErrGeneric                  = Error{500, "Something happened on our end, please contact the support", ""}
	API2ErrDatabase                 = Error{500, "Database Error", ""}
	API2ErrDuplicateUser            = Error{409, "User already taken", ""}
	API2ErrEmail                    = Error{500, "It was not possible to send the verification email", ""}
	API2ErrInvalidResponse          = Error{500, "Invalid response", ""}
	API2ErrTokenGeneration          = Error{500, "Token generation error", ""}
	API2ErrUUIDGeneration           = Error{500, "UUID generation failed", ""}
	API2ErrInvalidJSON              = Error{400, "Invalid JSON request", ""}
	API2ErrUnauthorized             = Error{401, "Unauthorized", ""}
	API2ErrInvalidCredentials       = Error{403, "Invalid Credentials", ""}
	API2ErrAccountDisabled          = Error{403, "Account disabled", ""}
	API2ErrInvalidURL               = Error{403, "URL could not be understood", ""}
	API2ErrNotFound                 = Error{404, "Resource not found", ""}
	API2ErrTokenExpired             = Error{404, "Token has expired", ""}
	API2ErrUserNotAdmin             = Error{422, "The specified user is not existing or is not an administrator", ""}
	API2ErrUserRegistrationDisabled = Error{409, "User registration is disabled", ""}
)

// RegistrationRequest contains all information to start the registtation
// process
type RegistrationRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
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

		if !s.config.EnableUserRegistration {
			jsonResponse(w, API2ErrUserRegistrationDisabled)
			return
		}

		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			jsonResponse(w, API2ErrInvalidJSON)
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

			if strings.Contains(err.Error(), "1062: Duplicate entry") {
				jsonResponse(w, API2ErrDuplicateUser)
				return
			} else {
				fmt.Println(err)
				jsonResponse(w, API2ErrGeneric)
				return
			}
		}

		uuid, err := uuid.NewV4()
		if err != nil {
			jsonResponse(w, API2ErrUUIDGeneration)
			return
		}

		var emailVerification = &EmailVerification{
			RegistrantID: registrant.ID,
			Email:        registrant.Email,
			Token:        uuid.String(),
			Timestamp:    time.Now().UTC().Unix(),
		}
		if err := s.store.CreateEmailVerification(emailVerification); err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		// build verification link for the email
		// TODO: make domain configurable
		var verificationLink = fmt.Sprintf(
			"http://%s%s/verify-email/%s",
			s.config.HTTPDomain,
			s.config.HTTPBasePath,
			emailVerification.Token,
		)

		if err := s.email.Email(
			emailVerification.Email,
			"K-Link-Registry: Please verify your email address",
			`html `+verificationLink,
			`hello, welcome to the K-Link registry. Please use this link to verify your mail address and set a password: `+verificationLink,
		); err != nil {
			jsonResponse(w, Error{422, err.Error(), ""})
			return
		}

		// success
		jsonResponse(w, response)
		return
	}
}

func (s *Server) handlePostPasswordReset() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}

func (s *Server) handleGetVerifyEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// get token from URL
		token := chi.URLParam(req, "token")

		// fetch EmailVerification
		verification, err := s.store.GetEmailVerificationByToken(token)
		if s.store.IsNotFound(err) {
			jsonResponse(w, API2ErrNotFound)
			return
		} else if err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		// terminate if the EmailVerification is invalid
		if verification.IsExpired() {
			jsonResponse(w, API2ErrTokenExpired)
			return
		}

		// fetch User
		user, err := s.store.GetRegistrantByID(verification.RegistrantID)
		if s.store.IsNotFound(err) {
			log.Printf("Verification: Registrant %d No longer exists", verification.RegistrantID)
			jsonResponse(w, API2ErrNotFound)
			return
		} else if err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		// return Response, based on the fact if this is the user's initial
		// email verification (after signup) the "require_password" parameter
		// will be set
		var response EmailVerificationModel
		response.DisplayName = user.Name
		response.RequirePassword = (len(user.Password) == 0)

		jsonResponse(w, response)
		return
	}
}

func (s *Server) handlePostVerifyEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// optional password, will be required when the user has no password
		// set yet, will not be used otherwise.
		var request SetPasswordRequest
		var response API2EmptyResponse

		// get token from URL
		token := chi.URLParam(req, "token")

		// fetch EmailVerification
		verification, err := s.store.GetEmailVerificationByToken(token)
		if s.store.IsNotFound(err) {
			jsonResponse(w, API2ErrNotFound)
			return
		} else if err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		// terminate if the EmailVerification is invalid
		if verification.IsExpired() {
			jsonResponse(w, API2ErrTokenExpired)
			return
		}

		// fetch User
		user, err := s.store.GetRegistrantByID(verification.RegistrantID)
		if s.store.IsNotFound(err) {
			jsonResponse(w, API2ErrNotFound)
			return
		} else if err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		// deserialize request
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			jsonResponse(w, API2ErrInvalidJSON)
			return
		}

		// set password if user has no password yet
		if len(user.Password) == 0 {
			user.SetPass(request.Password)
		}

		// change email to mail in EmailVerification
		user.Email = verification.Email

		// persist user
		if err := s.store.ReplaceRegistrant(user); err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}
		jsonResponse(w, response)
		return
	}
}

func (s *Server) handlePostSetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//// get token from URL
		//token := chi.URLParam(req, "token")
		//
		//// fetch PasswordReset
		//reset, err := s.store.GetPasswordResetByToken(token)
		//if err != nil {
		//	jsonResponse(w, API2ErrDatabase)
		//	return
		//}
		//
		//// terminate if PasswordReset is invalid
		//
		//// change user Password
	}
}
