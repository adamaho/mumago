package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"mumago/internal/api"
	"mumago/internal/db"
	"mumago/internal/realtime"
)

func main() {
	r := chi.NewRouter()

	db := db.New().Connect()
	todos := api.NewTodos()

	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		rt := realtime.New()
		todosApi := api.NewTodosApi(db, &rt, &todos)

		r.Get("/todos", todosApi.GetTodos)
		r.Post("/todos/{task}", todosApi.CreateTodo)
	})

	log.Println("starting todos server at https://localhost:3000")

	err := http.ListenAndServeTLS("localhost:3000", "cert.pem", "key.pem", r)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
