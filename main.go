package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/subscribe", subscribe)
	r.Post("/publish", publish)

	log.Println("starting server at https://localhost:3000")

	err := http.ListenAndServeTLS("localhost:3000", "cert.pem", "key.pem", r)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

}

func subscribe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("subscribe"))
}

func publish(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("publish"))
}
