package api

import (
	"encoding/json"
	"muma/internal/db"
	"muma/internal/helpers"
	"net/http"

	"gorm.io/gorm"
)

type TableApi struct {
	db *gorm.DB
}

// Creates a new TableApi instance
func NewTableApi(db *gorm.DB) TableApi {
	return TableApi{db: db}
}

type TableCreate struct {
	Name string
}

// Gets a new random table code
func (tApi *TableApi) CreateTable(w http.ResponseWriter, req *http.Request) {
	var t TableCreate

	err := json.NewDecoder(req.Body).Decode(&t)

	if err != nil {
		helpers.Log(helpers.Error, "Failed to parse CreateTable request body", err)
		helpers.HttpError(w, helpers.InvalidRequestBody, "")
		return
	}

	id, err := db.CreateTable(tApi.db, t.Name)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
		return
	}

	table, err := db.GetTableByID(tApi.db, id)

	if err != nil {
		helpers.HttpError(w, helpers.DatabaseError, "")
		return
	}

	tableJson, err := json.Marshal(table)

	if err != nil {
		helpers.HttpError(w, helpers.MarshalError, "")
		helpers.Log(helpers.Error, "Failed to marshal new table", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tableJson)

}
