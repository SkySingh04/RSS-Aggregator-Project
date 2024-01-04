package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("Hello World")

	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment ")
	}

	fmt.Println("PORT is ", portString)

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

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server is running on port %s", portString)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
