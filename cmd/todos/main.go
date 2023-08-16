package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"mumago/internal/realtime"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		rt := realtime.New()
		r.Get("/todos", func(w http.ResponseWriter, r *http.Request) { getTodos(w, r, &rt) })
		r.Post("/todos/{message}", func(w http.ResponseWriter, r *http.Request) { createTodo(w, r, &rt) })
	})

	log.Println("starting todos server at https://localhost:3000")

	err := http.ListenAndServeTLS("localhost:3000", "cert.pem", "key.pem", r)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

type Todo struct {
	TodoId int    `json:"todo_id"`
	Task   string `json:"task"`
}

func getTodos(w http.ResponseWriter, r *http.Request, rt *realtime.Realtime) {
	todo := Todo{TodoId: 1, Task: "Hello world"}
	d, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, "Failed to parse json", http.StatusInternalServerError)
		return
	}
	rt.Stream(w, r, d)
}

func createTodo(w http.ResponseWriter, r *http.Request, rt *realtime.Realtime) {
	msg := chi.URLParam(r, "message")
	for _, client := range rt.Clients {
		*client.Channel <- msg
	}
	fmt.Fprintf(w, "Message sent: %s", msg)
}
