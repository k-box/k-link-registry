package klinkregistry

import (
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// LoginRequest contains the credentials for a user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// generateToken generates a random token for use in one-time tokens, such as
// email verification and password reset
func generateToken() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic("Error while generating UUID Token: " + err.Error())
	}

	return uuid.String()
}

// handleGetSession provides an endpoint to check and refresh user sessions.
// if the user is not authenticated, an "unauthorized" error will be returned.
// if the user is authenticated, a refreshed token will be echoed back.
func (s *Server) handleGetSession() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var response SessionResponse

		u := s.sessions.GetUser(req)
		if u == nil {
			jsonResponse(w, API2ErrUnauthorized)
			return
		}

		token, err := s.sessions.GenerateToken(*u, time.Now().Add(1*time.Hour))
		if err != nil {
			jsonResponse(w, API2ErrTokenGeneration)
			return
		}

		response.Token = token
		response.UserID = u.ID
		response.Role = u.Role

		jsonResponse(w, response)
		return
	}
}

// handleCreateSession provides an endpoint to create sessions. The
// endpoint expects a valid username/password pair inside the request body.
// On correct authorization a session token will be returned.
func (s *Server) handleCreateSession() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var request LoginRequest
		var response SessionResponse

		json.NewDecoder(req.Body).Decode(&request)

		registrant, err := s.store.GetRegistrantByEmail(request.Email)
		if err != nil {
			// User not found
			jsonResponse(w, API2ErrInvalidCredentials)
			return
		}
		if err := registrant.CheckPass(request.Password); err != nil {
			// PW invalid
			jsonResponse(w, API2ErrInvalidCredentials)
			return
		}

		// If user is inactive, deny login anyways
		if !registrant.Active {
			jsonResponse(w, API2ErrAccountDisabled)
			return
		}

		// save LastLogin timestamp
		registrant.LastLogin = time.Now().UTC().Unix()
		err = s.store.ReplaceRegistrant(registrant)
		if err != nil {
			jsonResponse(w, err.Error())
		}

		// set user
		sessionUser := User{
			ID:          registrant.ID,
			DisplayName: registrant.Name,
			Role:        registrant.Role,
		}
		token, err := s.sessions.GenerateToken(
			sessionUser,
			time.Now().Add(15*time.Minute),
		)
		if err != nil {
			jsonResponse(w, API2ErrTokenGeneration)
			return
		}

		response.Token = token
		response.UserID = sessionUser.ID
		response.Role = sessionUser.Role
		jsonResponse(w, response)
		return
	}
}
