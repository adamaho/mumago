package db

import (
	"crypto/rand"
	"io"
	"muma/internal/helpers"

	"gorm.io/gorm"
)

const CODE_LENGTH = 6

type Table struct {
	ID   uint   `json:"table_id" gorm:"primaryKey"`
	Name string `json:"name"`
	Code string `json:"code" gorm:"primaryKey;autoIncrement:false" gorm:"type:varchar(6)"`
}

// creates a random table code
var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func createTableCode() (string, error) {
	b := make([]byte, CODE_LENGTH)
	n, err := io.ReadAtLeast(rand.Reader, b, CODE_LENGTH)
	if n != CODE_LENGTH {
		return "", err
	}

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	return string(b), nil
}

// Creates a new Table with a randomly generated sessionID
func CreateTable(db *gorm.DB, name string) (uint, error) {
	code, err := createTableCode()

	if err != nil {
		return 0, err
	}

	t := Table{Name: name, Code: code}
	result := db.Create(&t)

	if result.Error != nil {
		helpers.Log(helpers.Error, "Failed to CreateTable in database", result.Error)
		return 0, result.Error
	}

	return t.ID, nil
}

// Gets a Table by the provided ID
func GetTableByID(db *gorm.DB, tableID uint) (Table, error) {
	var table Table
	result := db.First(&table, tableID)

	if result.Error != nil {
		helpers.Log(helpers.Error, "Failed to GetTableByID from database", result.Error)
		return Table{}, result.Error
	}

	return table, nil
}
