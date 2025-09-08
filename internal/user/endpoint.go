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

type Response struct {
	Status int    `json:"status"`
	Data   any    `json:"data,omitempty"`
	Err    string `json:"error,omitempty"`
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
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "invalid request format"})
			return
		}

		if createRequest.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "first name is required"})
			return
		}

		if createRequest.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "last name is required"})
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
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: user})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: user})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.GetAll()

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: users})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateRequest UpdateRequest

		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "invalid request format"})
			return
		}

		if updateRequest.FirstName != nil && *updateRequest.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "first name is required"})
			return
		}

		if updateRequest.LastName != nil && *updateRequest.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "last name is required"})
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
			json.NewEncoder(w).Encode(Response{Status: 404, Err: "user does not exist"})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: "success"})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(Response{Status: 404, Err: "user does not exists"})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: "success"})
	}
}
