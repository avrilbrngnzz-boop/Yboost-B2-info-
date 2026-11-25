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
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Recipes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		ingredients TEXT NOT NULL,
		country_of_origin TEXT,
		dish_type TEXT, /*exemple: entree, plat, dessert, boisson, etc.*/
		diet_type TEXT, /*exemple: vegan, vegetarien, sans gluten, etc.*/
		instructions TEXT NOT NULL,
		photo_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS Comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		recipe_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE,
		FOREIGN KEY(recipe_id) REFERENCES Recipes(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Favorites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		recipe_id INTEGER NOT NULL,
		UNIQUE(user_id, recipe_id),
		FOREIGN KEY(user_id) REFERENCES Users(id) ON DELETE CASCADE,
		FOREIGN KEY(recipe_id) REFERENCES Recipes(id) ON DELETE CASCADE
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
	//Ajout de COALESCE pour transformer le NULL de SQLite en chaine vide ""
	err := db.QueryRow("SELECT COALESCE(photo_profil, '') FROM Users WHERE email = ?", email).Scan(&pp)
	if err != nil && err != sql.ErrNoRows {
		ErrorLogger.Printf("Erreur récupération photo_profil pour %s: %v", email, err)
	}
	return pp, err
}
