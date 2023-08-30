package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	todoID, err := db.CreateTodo(tApi.db, sessionID, task)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
		return
	}

	targetDb, err := db.GetTodosBySessionID(tApi.db, sessionID)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
		return
	}

	targetRealtime := realtime.Data{Data: targetDb}
	target, err := json.Marshal(targetRealtime)

	if err != nil {
		helpers.HttpError(w, helpers.MarshalError, "")
		helpers.Log(helpers.Error, "Failed to marshal target", err)
	}

	patch, err := tApi.rt.GeneratePatch(target, sessionID)

	if err != nil {
		helpers.HttpError(w, helpers.PatchError, "")
		helpers.Log(helpers.Error, "Failed to generate patch", err)
	}

	tApi.rt.PublishMsg(patch, sessionID)

	newTodo, err := db.GetTodoByID(tApi.db, todoID)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
	}

	newTodoJson, err := json.Marshal(newTodo)

	// TODO: figure out a way to not have to do this down the chain
	if err != nil {
		helpers.HttpError(w, helpers.MarshalError, "")
		helpers.Log(helpers.Error, "Failed to marshal new todos", err)
	}

	// TODO: create a helper function for this
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(newTodoJson)
}

// Updates a todo
func (tApi *TodosApi) UpdateTodo(w http.ResponseWriter, req *http.Request) {
	sessionID := chi.URLParam(req, "sessionID")
	taskIDParam := chi.URLParam(req, "taskID")

	taskID, e := strconv.ParseUint(taskIDParam, 10, 64)

	if e != nil {
		helpers.Log(helpers.Error, "Failed to parse taskID from params", e)
		helpers.HttpError(w, helpers.InvalidRequestParams, "")
		return
	}

	var t db.TodoForm
	err := json.NewDecoder(req.Body).Decode(&t)

	if err != nil {
		helpers.Log(helpers.Error, "Failed to parse TodoForm request body", err)
		helpers.HttpError(w, helpers.InvalidRequestBody, "")
		return
	}

	todoID, err := db.UpdateTodoByID(tApi.db, uint(taskID), sessionID, t)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
		return
	}

	// TODO: Figure out how to generalize this code so that it isnt so crazy. Pretty much the same thing from here down in order to create a new patch
	targetDb, err := db.GetTodosBySessionID(tApi.db, sessionID)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
		return
	}

	targetRealtime := realtime.Data{Data: targetDb}
	target, err := json.Marshal(targetRealtime)

	if err != nil {
		helpers.HttpError(w, helpers.MarshalError, "")
		helpers.Log(helpers.Error, "Failed to marshal target", err)
	}

	patch, err := tApi.rt.GeneratePatch(target, sessionID)

	if err != nil {
		helpers.HttpError(w, helpers.PatchError, "")
		helpers.Log(helpers.Error, "Failed to generate patch", err)
	}

	tApi.rt.PublishMsg(patch, sessionID)

	updatedTodo, err := db.GetTodoByID(tApi.db, todoID)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
	}

	updatedTodoJson, err := json.Marshal(updatedTodo)

	// TODO: figure out a way to not have to do this down the chain
	if err != nil {
		helpers.HttpError(w, helpers.MarshalError, "")
		helpers.Log(helpers.Error, "Failed to marshal updated todo", err)
	}

	// TODO: create a helper function for this
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(updatedTodoJson)
}
