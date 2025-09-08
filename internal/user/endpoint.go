package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

type UpdateRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
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

		user, err := s.Create(
			createRequest.FirstName,
			createRequest.LastName,
			createRequest.Email,
			createRequest.Phone,
		)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.GetAll()

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateRequest UpdateRequest

		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"invalid request format"})
			return
		}

		if updateRequest.FirstName != nil && *updateRequest.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"first name is required"})
			return
		}

		if updateRequest.LastName != nil && *updateRequest.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"last name is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		err := s.Update(
			id,
			updateRequest.FirstName,
			updateRequest.LastName,
			updateRequest.Email,
			updateRequest.Phone,
		)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorResponse{"user does not exist"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"data": "success"})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorResponse{"user does not exists"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"data": "success"})
	}
}
