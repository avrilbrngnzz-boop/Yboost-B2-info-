package yboost

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("DB non initialisée")
	}
	return db, nil
}

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "yboost.db")
	if err != nil {
		panic(err)
	}
}

func CreateDB() {
	InitDB()

	schema := `
	CREATE TABLE IF NOT EXISTS Users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		photo_profil TEXT 
	);

	CREATE TABLE IF NOT EXISTS Session (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		token TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		fmt.Println("Erreur lors de la création des tables:", err)
		panic(err)
	}
}

func getUserPP(email string) (string, error) {
	var pp string
	err := db.QueryRow("SELECT photo_profil FROM Users WHERE email = ?", email).Scan(&pp)
	return pp, err
}
