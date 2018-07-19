package klinkregistry

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// Possible values for roles of the registrant
const (
	RoleUser  = "ROLE_USER"
	RoleAdmin = "ROLE_ADMIN"
	RoleOwner = "ROLE_OWNER"
)

// RegistrantModel is the JSON representation of a registrant
type RegistrantModel struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Password  []byte `json:"-"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
	LastLogin int64  `json:"last_login"`
}

// SetPasswordRequest contains a new password to be set
type SetPasswordRequest struct {
	Password string `json:"password"`
}

// handleListRegistrants provides an endpoint that returns a list of all
// registrants inside the database
func (s *Server) handleListRegistrants() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var responses []RegistrantModel

		registrants, err := s.store.ListRegistrants()
		if s.store.IsNotFound(err) {
			jsonResponse(w, responses)
			return
		} else if err != nil {
			log.Println(err.Error())
			jsonResponse(w, API2ErrDatabase)
			return
		}

		user := s.sessions.GetUser(req)

		for _, registrant := range registrants {
			// do not list registrants  if the registrant is an USER
			// (not ADMIN or OWNER)
			if user.Role == RoleUser && registrant.ID != user.ID {
				continue
			}

			responses = append(responses, RegistrantModel(*registrant))
		}

		jsonResponse(w, responses)
	}
}

func (s *Server) handleCreateRegistrant() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var request RegistrantModel
		var response RegistrantModel

		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			jsonResponse(w, API2ErrInvalidJSON)
			return
		}

		user := s.sessions.GetUser(req)

		if user.Role != RoleAdmin && user.Role != RoleOwner {
			jsonResponse(w, API2ErrUnauthorized)
			return
		}

		registrant := Registrant(request)

		if err := s.store.CreateRegistrant(&registrant); err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		response = RegistrantModel(registrant)

		jsonResponse(w, response)
		return
	}
}

// handleGetRegistrant provides an endpoint that returns a single registrant
// from the database
func (s *Server) handleGetRegistrant() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var response RegistrantModel

		idString := chi.URLParam(req, "id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			jsonResponse(w, API2ErrInvalidURL)
			return
		}

		registrant, err := s.store.GetRegistrantByID(id)
		if s.store.IsNotFound(err) {
			jsonResponse(w, API2ErrNotFound)
			return
		} else if err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		user := s.sessions.GetUser(req)
		// do not show unowned applications if the registrant is an USER
		// (not ADMIN or OWNER)
		if user.Role == RoleUser && registrant.ID != user.ID {
			jsonResponse(w, API2ErrUnauthorized)
			return
		}

		response = RegistrantModel(*registrant)
		jsonResponse(w, response)
		return
	}

}

// handleUpdateRegistrant provides an endpoint that allows updating of
// attributes for registrants
func (s *Server) handleUpdateRegistrant() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var response RegistrantModel
		var request RegistrantModel

		idString := chi.URLParam(req, "id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			jsonResponse(w, API2ErrInvalidURL)
			return
		}

		registrant, err := s.store.GetRegistrantByID(id)
		if s.store.IsNotFound(err) {
			jsonResponse(w, API2ErrNotFound)
			return
		} else if err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		user := s.sessions.GetUser(req)

		// check if we have permission to edit this registrants
		if user.ID == registrant.ID {
			// we must either be the registrant itself ..
		} else if user.Role == RoleOwner || user.Role == RoleAdmin {
			// .. or have a role of admin or owner.
		} else {
			// otherwise we can strop processing here
			jsonResponse(w, API2ErrInvalidCredentials)
			return
		}

		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			jsonResponse(w, API2ErrInvalidJSON)
			return
		}

		// use the registrant as a base to apply our request to:
		// registrant.ID must stay the same.
		registrant.Name = request.Name
		registrant.Email = request.Email

		// allow change of some attributes, if user is admin or owner
		if user.Role == RoleOwner || user.Role == RoleAdmin {
			registrant.Active = request.Active
			registrant.Email = request.Email
			registrant.Role = request.Role
		}

		if err := s.store.ReplaceRegistrant(registrant); err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		response = RegistrantModel(*registrant)
		jsonResponse(w, response)
		return
	}
}

// handleDeleteRegistrant provides an endpoint that allows deletion of
// registrants
func (s *Server) handleDeleteRegistrant() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := chi.URLParam(req, "id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			jsonResponse(w, API2ErrInvalidURL)
			return
		}

		registrant, err := s.store.GetRegistrantByID(id)
		if s.store.IsNotFound(err) {
			// deletion of already deleted entry should succeed
			jsonResponse(w, API2EmptyResponse{})
			return
		} else if err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		user := s.sessions.GetUser(req)

		if registrant.Role == RoleUser || user.Role != RoleUser {
			// allow admins to remove simple users
		} else if user.Role == "ROLE_OWNER" {
			// an owner role may remove everything
		} else {
			jsonResponse(w, API2ErrUnauthorized)
			return
		}

		if err := s.store.DeleteRegistrant(registrant.ID); err != nil {
			jsonResponse(w, API2ErrDatabase)
			return
		}

		jsonResponse(w, API2EmptyResponse{})
		return
	}
}
