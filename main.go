package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/S3ergio31/curso-go-seccion-4/internal/course"
	"github.com/S3ergio31/curso-go-seccion-4/internal/user"
	"github.com/S3ergio31/curso-go-seccion-4/pkg/bootstrap"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db, err := bootstrap.DBConnection()
	logger := bootstrap.InitLogger()

	if err != nil {
		logger.Fatalln(err)
	}

	router := mux.NewRouter()

	userRepository := user.NewRepository(logger, db)
	userService := user.NewService(userRepository, logger)
	userEndpoints := user.MakeEndpoints(userService)

	router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoints.Get).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoints.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEndpoints.Delete).Methods("DELETE")

	courseRepository := course.NewRepository(logger, db)
	courseService := course.NewService(courseRepository, logger)
	courseEndpoints := course.MakeEndpoints(courseService)

	router.HandleFunc("/courses", courseEndpoints.Create).Methods("POST")
	router.HandleFunc("/courses", courseEndpoints.GetAll).Methods("GET")
	router.HandleFunc("/courses/{id}", courseEndpoints.Get).Methods("GET")
	router.HandleFunc("/courses/{id}", courseEndpoints.Update).Methods("PATCH")
	router.HandleFunc("/courses/{id}", courseEndpoints.Delete).Methods("DELETE")

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("APP_URL"), os.Getenv("APP_PORT")),
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}

	logger.Fatal(server.ListenAndServe())
}
