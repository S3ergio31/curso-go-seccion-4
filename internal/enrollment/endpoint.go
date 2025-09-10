package enrollment

import (
	"encoding/json"
	"net/http"

	"github.com/S3ergio31/curso-go-seccion-4/pkg/meta"
)

type Controller func(w http.ResponseWriter, r *http.Request)

type Endpoints struct {
	Create Controller
}

type CreateRequest struct {
	UserID   string `json:"user_id"`
	CourseID string `json:"course_id"`
}

type Response struct {
	Status int        `json:"status"`
	Data   any        `json:"data,omitempty"`
	Err    string     `json:"error,omitempty"`
	Meta   *meta.Meta `json:"meta,omitempty"`
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
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

		if createRequest.UserID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "name is required"})
			return
		}

		if createRequest.CourseID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "start date is required"})
			return
		}

		enrollment, err := s.Create(
			createRequest.UserID,
			createRequest.CourseID,
		)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: enrollment})
	}
}
