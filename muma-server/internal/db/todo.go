package db

import (
	"log"
	"muma/internal/helpers"

	"gorm.io/gorm"
)

type Todo struct {
	ID        uint   `json:"todo_id" gorm:"primaryKey"`
	Task      string `json:"task"`
	Checked   bool   `json:"checked" gorm:"default:false"`
	SessionID string `json:"session_id"`
}

// Creates a new todo
func CreateTodo(db *gorm.DB, sessionID string, task string) (uint, error) {
	t := Todo{Task: task, SessionID: sessionID}
	result := db.Create(&t)

	if result.Error != nil {
		helpers.Log(helpers.Error, "Failed to CreateTodo in database", result.Error)
		return 0, result.Error
	}

	return t.ID, nil
}

// Gets all todos
func GetTodos(db *gorm.DB) ([]Todo, error) {
	var todos []Todo
	result := db.Find(&todos)

	if result.Error != nil {
		log.Print("Failed to get todos:", result.Error)
		return nil, result.Error
	}

	return todos, nil
}

// Gets all todos for a specific sessionID
func GetTodosBySessionID(db *gorm.DB, sessionID string) ([]Todo, error) {
	var todos []Todo
	result := db.Where("session_id = ?", sessionID).Find(&todos)

	if result.Error != nil {
		helpers.Log(helpers.Error, "Failed to GetTodosBySessionID from database", result.Error)
		return nil, result.Error
	}

	return todos, nil
}

// Gets a single todo by ID
func GetTodoByID(db *gorm.DB, todoID uint) (Todo, error) {
	var todo Todo
	result := db.First(&todo, todoID)

	if result.Error != nil {
		helpers.Log(helpers.Error, "Failed to GetTodoByID from database", result.Error)
		return Todo{}, result.Error
	}

	return todo, nil
}

type TodoForm struct {
	Task    string `json:"task"`
	Checked bool   `json:"checked"`
}

// Gets all todos for a specific sessionID
func UpdateTodoByID(db *gorm.DB, todoID uint, sessionID string, todo TodoForm) (uint, error) {
	result := db.Save(&Todo{ID: todoID, SessionID: sessionID, Task: todo.Task, Checked: todo.Checked})

	if result.Error != nil {
		helpers.Log(helpers.Error, "Failed to UpdateTodoByID from database", result.Error)
		return 0, result.Error
	}

	return todoID, nil
}
