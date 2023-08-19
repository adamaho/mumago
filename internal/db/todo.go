package db

import (
	"log"

	"gorm.io/gorm"
)

type Todo struct {
	ID   uint   `json:"todo_id" gorm:"primaryKey"`
	Task string `json:"task"`
}

// Creates a new todo
func CreateTodo(db *gorm.DB, task string) (uint, error) {
	t := Todo{Task: task}
	result := db.Create(&t)

	if result.Error != nil {
		log.Print("Failed to create new todo:", result.Error)
		return 0, result.Error
	}

	return t.ID, nil
}

// Gets all todos
func GetTodos(db *gorm.DB) ([]Todo, error) {
	var todos []Todo
	result := db.Find(&todos)

	if result.Error != nil {
		log.Print("Failed to create new todo:", result.Error)
		return nil, result.Error
	}

	return todos, nil
}

// Gets a single todo by ID
func GetTodoByID(db *gorm.DB, todoID uint) (Todo, error) {
	var todo Todo
	result := db.First(&todo, todoID)

	if result.Error != nil {
		log.Print("Failed to get todo with id:", todoID)
		return Todo{}, result.Error
	}

	return todo, nil
}
