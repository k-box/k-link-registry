package klinkregistry

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

// EmailVerificationRequest contains information for an email verification
type EmailVerificationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// handleGetVerifyEmail is called by the frontend to get additional
// information about an email verification.
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
		if len(user.Password) == 0 {
			// If no password is set, always require user to set a password
			response.RequirePassword = true
		} else if verification.Email == user.Email {
			// If the email is unchanged, this is a password reset request
			response.RequirePassword = true
		}

		jsonResponse(w, response)
		return
	}
}

// handlePostVerifyEmail is called by the frontend to perform a email
// verification for a existing token
func (s *Server) handlePostVerifyEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// optional password, will be required when the user has no password
		// set yet, will not be used otherwise.
		var request EmailVerificationRequest
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

		if len(user.Password) == 0 {
			// set password if user has no password yet
			if len(request.Password) == 0 {
				jsonResponse(w, API2ErrPasswordRequired)
				return
			}
			user.SetPass(request.Password)
		} else if user.Email == request.Email {
			// If email is unchanged, this is a password reset token, and we
			// update the password
			if len(request.Password) == 0 {
				jsonResponse(w, API2ErrPasswordRequired)
				return
			}
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

// CreateVerification creates an email verification token, sends it and stores
// it in the Database.
// * If User Password is empty, a "signup confirmation" email will be sent.
// * If newmail is not empty and differs from the current user email,
//   a "email confirmation" Email will be sent.
// * Otherwise, a "password reset" email will be dispatched.
func (s *Server) CreateVerification(registrant *Registrant, newmail string) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return API2ErrUUIDGeneration
	}

	// First, delete any verifications that might still exist.
	// Error handling is ommited, since we do not need to recover from failure
	// here.
	_ = s.store.DeleteEmailVerification(registrant.ID)

	if len(newmail) == 0 {
		// if the new mail is not set, use the registrant's current mail
		newmail = registrant.Email
	}

	var emailVerification = &EmailVerification{
		RegistrantID: registrant.ID,
		Email:        newmail,
		Token:        uuid.String(),
		Timestamp:    time.Now().UTC().Unix(),
	}
	if err := s.store.CreateEmailVerification(emailVerification); err != nil {
		return API2ErrDatabase
	}

	// build verification link for the email
	// TODO: make domain configurable
	var verificationLink = fmt.Sprintf(
		"https://%s%s/auth/confirm/%s",
		s.config.HTTPDomain,
		s.config.HTTPBasePath,
		emailVerification.Token,
	)

	// TODO: localize somehow?
	if err := s.email.Email(
		emailVerification.Email,
		"K-Link-Registry: Please verify your email address",
		`html `+verificationLink,
		`hello, welcome to the K-Link registry. Please use this link to verify your mail address and set a password: `+verificationLink,
	); err != nil {
		return API2ErrEmail
	}
	return nil

}
