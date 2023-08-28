package api

import (
	"encoding/json"
	"net/http"

	"muma/internal/db"
	"muma/internal/helpers"
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

	todosData, err := db.GetTodosBySessionID(tApi.db, sessionID)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
		return
	}

	todos := realtime.Data{Data: todosData}
	todosJson, err := json.Marshal(todos)

	if err != nil {
		helpers.HttpError(w, helpers.MarshalError, "")
		helpers.Log(helpers.Error, "Failed to marshal todos", err)
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
		helpers.HttpError(w, helpers.DatabaseError, "")
		return
	}

	// fetch all of the todos from the db
	targetDb, err := db.GetTodosBySessionID(tApi.db, sessionID)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
		return
	}

	// marshal todos to json
	targetRealtime := realtime.Data{Data: targetDb}
	target, err := json.Marshal(targetRealtime)

	if err != nil {
		helpers.HttpError(w, helpers.MarshalError, "")
		helpers.Log(helpers.Error, "Failed to marshal target", err)
	}

	tApi.rt.PublishPatch(target, sessionID)

	// fetch the new todo from the db to return to the user
	newTodo, err := db.GetTodoByID(tApi.db, todoID)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
	}

	newTodoJson, err := json.Marshal(newTodo)

	if err != nil {
		helpers.HttpError(w, helpers.MarshalError, "")
		helpers.Log(helpers.Error, "Failed to marshal new todos", err)
	}

	// TODO: create a helper function for this
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(newTodoJson)
}
