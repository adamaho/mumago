package db

type Todo struct {
	ID   uint   `json:"todo_id" gorm:"primaryKey"`
	Task string `json:"task"`
}
