package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func OpenDB() error {
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func CloseDB() error {
	return DB.Close()
}

func SetupDB() error {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT,
		created_at TEXT,
		due_date TEXT,
		completed BOOLEAN);`)
	if err != nil {
		return err
	}
	return nil
}
