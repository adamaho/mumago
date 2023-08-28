package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	dsn string
}

// Creates a new database instance
func New() *Database {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Failed to get env from .env")
	}

	return &Database{
		dsn: os.Getenv("DATABASE_DSN"),
	}
}

// Creates a new connection to the database
func (d *Database) Connect() *gorm.DB {
	db, err := gorm.Open(mysql.Open(d.dsn), &gorm.Config{})
	db.AutoMigrate(&Todo{}, &Table{})

	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	return db
}
