package api

import (
	"encoding/json"
	"net/http"

	"muma/internal/db"
	"muma/internal/realtime"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type TodosApi struct {
	db *gorm.DB
	rt *realtime.Realtime
}

// Creates a new TodosApi instance
func NewTodosApi(db *gorm.DB) TodosApi {
	rt := realtime.New()
	return TodosApi{db: db, rt: &rt}
}

// Returns all todos or an optional stream of todos
func (tApi *TodosApi) GetTodos(w http.ResponseWriter, req *http.Request) {
	sessionID := chi.URLParam(req, "sessionID")

	todosData, err := db.GetTodos(tApi.db)

	if err != nil {
		http.Error(w, "Failed to get todos from db", http.StatusInternalServerError)
		return
	}

	todos := realtime.Data{Data: todosData}
	todosJson, err := json.Marshal(todos)

	if err != nil {
		http.Error(w, "Failed to marshal todos to json", http.StatusInternalServerError)
		return
	}

	tApi.rt.Stream(w, req, todosJson, sessionID)
}

// Creates a new todo
func (tApi *TodosApi) CreateTodo(w http.ResponseWriter, req *http.Request) {
	sessionID := chi.URLParam(req, "sessionID")
	task := chi.URLParam(req, "task")

	// create the todo
	todoID, err := db.CreateTodo(tApi.db, sessionID, task)

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
	targetRealtime := realtime.Data{Data: targetDb}
	target, err := json.Marshal(targetRealtime)

	if err != nil {
		http.Error(w, "Failed to marshal new todos to json", http.StatusInternalServerError)
	}

	tApi.rt.PublishPatch(target, sessionID)

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
