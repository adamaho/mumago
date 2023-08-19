package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"mumago/internal/db"
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

	db.New().Connect()

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
	// TODO: get todos from db instead
	d, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, "Failed to parse json", http.StatusInternalServerError)
		return
	}
	rt.Stream(w, r, d)
}

func createTodo(w http.ResponseWriter, r *http.Request, rt *realtime.Realtime, todos *Todos) {
	task := chi.URLParam(r, "task")
	ts := todos.AddTodo(Todo{TodoId: uuid.New(), Task: task})
	target, _ := json.Marshal(ts)

	// TODO: handle error

	rt.PublishPatch(w, target)
}
