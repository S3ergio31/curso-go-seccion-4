package course

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/S3ergio31/curso-go-seccion-4/pkg/meta"
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
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type UpdateRequest struct {
	Name      *string `json:"name"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

type Response struct {
	Status int        `json:"status"`
	Data   any        `json:"data,omitempty"`
	Err    string     `json:"error,omitempty"`
	Meta   *meta.Meta `json:"meta"`
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

		if createRequest.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "name is required"})
			return
		}

		if createRequest.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "start date is required"})
			return
		}

		if createRequest.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "end date is required"})
			return
		}

		course, err := s.Create(
			createRequest.Name,
			createRequest.StartDate,
			createRequest.EndDate,
		)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: course})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		course, err := s.Get(id)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: course})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filters := Filters{
			Name:      query.Get("name"),
			StartDate: query.Get("start_date"),
			EndDate:   query.Get("end_date"),
		}

		count, err := s.Count(filters)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(Response{Status: 500, Err: err.Error()})
			return
		}

		limit, _ := strconv.Atoi(query.Get("limit"))
		page, _ := strconv.Atoi(query.Get("page"))
		meta, err := meta.New(page, limit, count)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(Response{Status: 500, Err: err.Error()})
			return
		}

		courses, err := s.GetAll(filters, meta.Offset(), meta.Limit())

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{
			Status: 200,
			Data:   courses,
			Meta:   meta,
		})
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

		if updateRequest.Name != nil && *updateRequest.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "name is required"})
			return
		}

		if updateRequest.StartDate != nil && *updateRequest.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "start date is required"})
			return
		}

		if updateRequest.EndDate != nil && *updateRequest.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "end date is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		err := s.Update(
			id,
			updateRequest.Name,
			updateRequest.StartDate,
			updateRequest.EndDate,
		)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(Response{Status: 404, Err: "course does not exist"})
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
			json.NewEncoder(w).Encode(Response{Status: 404, Err: "course does not exists"})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: "success"})
	}
}
