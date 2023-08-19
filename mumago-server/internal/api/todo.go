package api

import (
	"encoding/json"
	"mumago/internal/db"
	"mumago/internal/realtime"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Todos struct {
	Data []db.Todo `json:"data"`
}

type TodosApi struct {
	db *gorm.DB
	rt *realtime.Realtime
}

// Creates a new TodosApi instance
func NewTodosApi(db *gorm.DB, rt *realtime.Realtime) TodosApi {
	return TodosApi{db, rt}
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

	// fetch all of the todos from the db
	targetDb, err := db.GetTodos(tApi.db)

	if err != nil {
		http.Error(w, "Failed to get new todos from db", http.StatusInternalServerError)
		return
	}

	// marshal todos to json
	targetStruct := Todos{Data: targetDb}
	target, err := json.Marshal(targetStruct)

	if err != nil {
		http.Error(w, "Failed to marshal new todos to json", http.StatusInternalServerError)
	}

	tApi.rt.PublishPatch(target)

	// fetch the new todo from the db to return to the user
	newTodo, err := db.GetTodoByID(tApi.db, todoID)

	if err != nil {
		http.Error(w, "Failed to get new todo from db", http.StatusInternalServerError)
	}

	newTodoJson, err := json.Marshal(newTodo)

	if err != nil {
		http.Error(w, "Failed to marshal new todo to json", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(newTodoJson)
}
