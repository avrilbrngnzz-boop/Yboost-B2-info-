package yboost

import (
	"database/sql"
	"fmt"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDB() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("DB non initialisée")
	}
	return db, nil
}

func InitDB() {
	// var err error
	// db, err = sql.Open("sqlite3", "yboost.db")
	// if err != nil {
	// 	panic(err)
	// }

	//SCRUM-63 Verification explicite de l'existence du fichier .db
	const dbPath = "yboost.db"
	_, err := os.Stat(dbPath)
	dbExists := !errors.Is(err, os.ErrNotExist)
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		ErrorLogger.Fatalf("Erreur de l'ouverture SQLite: %v", err)
	}
	//Verification que la connexion répond réellement
	if err = db.Ping(); err != nil {
		ErrorLogger.Fatalf("Impossible de joindre la DB: %v", err)
	}
	//Activation du ON DELETE CASCADE defini dans la table de Session
	if _, err = db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		ErrorLogger.Fatalf("Impossible d'activer les foreign_keys: %v", err)
	}
	if !dbExists {
		InfoLogger.Printf("Fichier %s introuvable: creation automatique effecuée.", dbPath)
	} else {
		InfoLogger.Printf("Connexion reussie a la base de donnees %s.", dbPath)
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

	CREATE TABLE IF NOT EXISTS Posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS Comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES Posts(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		UNIQUE(user_id, post_id),
		FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES Posts(id) ON DELETE CASCADE
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		//fmt.Println("Erreur lors de la création des tables:", err)
		//panic(err)
		ErrorLogger.Fatalf("Erreur lors de la création des tables: %v", err)
	}
	InfoLogger.Println("Verification/creation des tables terminee.")
}

func getUserPP(email string) (string, error) {
	var pp string
	err := db.QueryRow("SELECT photo_profil FROM Users WHERE email = ?", email).Scan(&pp)
	if err != nil && err != sql.ErrNoRows {
		ErrorLogger.Printf("Erreur récupération photo_profil pour %s: %v", email, err)
	}
	return pp, err
}
