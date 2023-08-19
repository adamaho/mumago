package api

import (
	"encoding/json"
	"mumago/internal/realtime"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	TodoId uuid.UUID `json:"todo_id"`
	Task   string    `json:"task"`
}

type Todos struct {
	Data []Todo `json:"data"`
}

// Creates a new Todos struct
func NewTodos() Todos {
	return Todos{Data: make([]Todo, 0)}
}

// Adds a todo to the database
func (t *Todos) AddTodo(todo Todo) Todos {
	t.Data = append(t.Data, todo)
	return *t
}

type TodosApi struct {
	db    *gorm.DB
	rt    *realtime.Realtime
	todos *Todos
}

// Creates a new TodosApi instance
func NewTodosApi(db *gorm.DB, rt *realtime.Realtime, todos *Todos) TodosApi {
	return TodosApi{db, rt, todos}
}

// Returns all todos or an optional stream of todos
func (t *TodosApi) GetTodos(w http.ResponseWriter, req *http.Request) {
	// TODO: Get todos from db
	d, err := json.Marshal(t.todos)
	if err != nil {
		http.Error(w, "Failed to parse json", http.StatusInternalServerError)
		return
	}
	t.rt.Stream(w, req, d)
}

// Creates a new todo
func (t *TodosApi) CreateTodo(w http.ResponseWriter, req *http.Request) {
	// TODO: create todo in db
	task := chi.URLParam(req, "task")
	ts := t.todos.AddTodo(Todo{TodoId: uuid.New(), Task: task})
	target, _ := json.Marshal(ts)

	// TODO: handle error

	t.rt.PublishPatch(w, target)
}
