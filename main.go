package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/S3ergio31/curso-go-seccion-4/internal/user"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	db := createDatabase()
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	router := mux.NewRouter()

	userRepository := user.NewRepository(logger, db)
	userService := user.NewService(userRepository, logger)
	userEndpoints := user.MakeEndpoints(userService)

	router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc("/users", userEndpoints.Update).Methods("PATCH")
	router.HandleFunc("/users", userEndpoints.Delete).Methods("DELETE")

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("APP_URL"), os.Getenv("APP_PORT")),
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}

	log.Fatal(server.ListenAndServe())
}

func createDatabase() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	log.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db = db.Debug()
	db.AutoMigrate(&user.User{})

	return db
}
