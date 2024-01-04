package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Akash-Singh04/rssaggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Hello World")

	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment ")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("PORT is not found in the environment ")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	c := (cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		MaxAge:           300,
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false}))
	router.Use(c.Handler)

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness) //Only works on GET requests
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server is running on port %s", portString)
	error := server.ListenAndServe()
	if error != nil {
		log.Fatal(err)
	}
}
