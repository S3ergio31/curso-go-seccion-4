package user

import (
	"encoding/json"
	"net/http"
)

type Controller func(w http.ResponseWriter, r *http.Request)

type Endpoints struct {
	Create Controller
	Get    Controller
	GetAll Controller
	Update Controller
	Delete Controller
}

type CreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var createRequest CreateRequest

		if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"invalid request format"})
			return
		}

		if createRequest.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"first name is required"})
			return
		}

		if createRequest.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"last name is required"})
			return
		}
		s.Create(
			createRequest.FirstName,
			createRequest.LastName,
			createRequest.Email,
			createRequest.Phone,
		)
		json.NewEncoder(w).Encode(createRequest)
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}
