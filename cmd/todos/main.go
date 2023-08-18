package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/mattbaird/jsonpatch"

	"mumago/internal/realtime"
)

type Todo struct {
	TodoId uuid.UUID `json:"todo_id"`
	Task   string    `json:"task"`
}

type Todos struct {
	Data []Todo `json:"data"`
}

// Adds a todo to the database
func (t *Todos) AddTodo(todo Todo) Todos {
	t.Data = append(t.Data, todo)
	return *t
}

func main() {
	r := chi.NewRouter()

	todos := Todos{Data: make([]Todo, 0)}

	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		rt := realtime.New()
		r.Get("/todos", func(w http.ResponseWriter, r *http.Request) { getTodos(w, r, &rt, &todos) })
		r.Post("/todos/{task}", func(w http.ResponseWriter, r *http.Request) { createTodo(w, r, &rt, &todos) })
	})

	log.Println("starting todos server at https://localhost:3000")

	err := http.ListenAndServeTLS("localhost:3000", "cert.pem", "key.pem", r)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request, rt *realtime.Realtime, todos *Todos) {
	// TODO: connect to a db
	d, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, "Failed to parse json", http.StatusInternalServerError)
		return
	}
	rt.Stream(w, r, d)
}

func createTodo(w http.ResponseWriter, r *http.Request, rt *realtime.Realtime, todos *Todos) {
	// TODO: connect to a db
	task := chi.URLParam(r, "task")
	t := Todo{TodoId: uuid.New(), Task: task}

	originalJson, _ := json.Marshal(*todos)

	target := todos.AddTodo(t)
	targetJson, _ := json.Marshal(target)

	patch, _ := jsonpatch.CreatePatch(originalJson, targetJson)
	patchJson, err := json.Marshal(patch)

	if err != nil {
		fmt.Println("err", err)
		http.Error(w, "Failed to marshal json", http.StatusInternalServerError)
		return
	}

	for _, client := range rt.Clients {
		*client.Channel <- patchJson
	}
}
