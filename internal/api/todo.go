package api

import (
	"encoding/json"
	"mumago/internal/db"
	"mumago/internal/realtime"
	"net/http"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type Todos struct {
	Data []db.Todo `json:"data"`
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
func (tApi *TodosApi) GetTodos(w http.ResponseWriter, req *http.Request) {
	todosData, err := db.GetTodos(tApi.db)

	// TODO: handle error here

	todos := Todos{Data: todosData}

	todosJson, err := json.Marshal(todos)

	// TODO: figure out a common way to handle json and json errors

	if err != nil {
		http.Error(w, "Failed to get todos from db", http.StatusInternalServerError)
		return
	}

	tApi.rt.Stream(w, req, todosJson)
}

// Creates a new todo
func (tApi *TodosApi) CreateTodo(w http.ResponseWriter, req *http.Request) {
	task := chi.URLParam(req, "task")

	// create the todo
	todoID, err := db.CreateTodo(tApi.db, task)

	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	// fetch the new todo from the db to return to the user

	// fetch all of the todos from the db
	targetData, err := db.GetTodos(tApi.db)

	if err != nil {
		http.Error(w, "Failed to get new todos from db", http.StatusInternalServerError)
		return
	}

	// marshal todos to json
	target, err := json.Marshal(targetData)

	if err != nil {
		http.Error(w, "Failed to marshal new todos to json", http.StatusInternalServerError)
	}

	tApi.rt.PublishPatch(w, target)
}
