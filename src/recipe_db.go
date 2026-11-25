package yboost

import (
	"database/sql"
	"time"
)

type Recipe struct {
	ID              int
	UserID          int
	Author          string
	AuthorPP        string
	Title           string
	Description     string
	Ingredients     string
	CountryOfOrigin string
	DishType        string
	DietType        string
	Instructions    string
	PhotoURL        string
	CreatedAt       time.Time
	FavoritesCount  int
}

type Comment struct {
	ID        int
	RecipeID  int
	UserID    int
	Author    string
	AuthorPP  string
	Content   string
	CreatedAt time.Time
}

func CreateRecipe(userID int, title, description, ingredients, country, dishType, dietType, instructions, photoURL string) error {
	query := `
		INSERT INTO Recipes (user_id, title, description, ingredients, country_of_origin, dish_type, diet_type, instructions, photo_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query, userID, title, description, ingredients, country, dishType, dietType, instructions, photoURL)
	if err != nil {
		ErrorLogger.Printf("Erreur création recette (user %d): %v", userID, err)
	}
	return err
}

func GetAllRecipes() ([]Recipe, error) {
	query := `
		SELECT Recipes.id, Recipes.user_id, Users.username, COALESCE(Users.photo_profil, ''), 
			Recipes.title, COALESCE(Recipes.description, ''), Recipes.ingredients, COALESCE(Recipes.country_of_origin, ''), 
			COALESCE(Recipes.dish_type, ''), COALESCE(Recipes.diet_type, ''), Recipes.instructions, COALESCE(Recipes.photo_url, ''), Recipes.created_at,
			(SELECT COUNT(*) FROM Favorites WHERE Favorites.recipe_id = Recipes.id) AS favorites_count
		FROM Recipes
		INNER JOIN Users ON Recipes.user_id = Users.id
		ORDER BY Recipes.created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		ErrorLogger.Printf("Erreur récupération recettes: %v", err)
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err := rows.Scan(
			&r.ID, &r.UserID, &r.Author, &r.AuthorPP, 
			&r.Title, &r.Description, &r.Ingredients, 
			&r.CountryOfOrigin, &r.DishType, &r.DietType, 
			&r.Instructions, &r.PhotoURL, &r.CreatedAt, &r.FavoritesCount,
		)
		if err != nil {
			ErrorLogger.Printf("Erreur lecture ligne recette: %v", err)
			return nil, err
		}
		recipes = append(recipes, r)
	}
	return recipes, nil
}

func CreateComment(userID, recipeID int, content string) error {
	_, err := db.Exec("INSERT INTO Comments (user_id, recipe_id, content) VALUES (?, ?, ?)", userID, recipeID, content)
	if err != nil {
		ErrorLogger.Printf("Erreur création du commentaire (user %d, recipe %d): %v", userID, recipeID, err)
	}
	return err
}

func GetCommentsByRecipeID(recipeID int) ([]Comment, error) {
	query := `
		SELECT 
			Comments.id, Comments.recipe_id, Comments.user_id, 
			Users.username, COALESCE(Users.photo_profil, ''), 
			Comments.content, Comments.created_at
		FROM Comments
		INNER JOIN Users ON Comments.user_id = Users.id
		WHERE Comments.recipe_id = ?
		ORDER BY Comments.created_at ASC
	`
	rows, err := db.Query(query, recipeID)
	if err != nil {
		ErrorLogger.Printf("Erreur récupération commentaires (recipe %d): %v", recipeID, err)
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.RecipeID, &c.UserID, &c.Author, &c.AuthorPP, &c.Content, &c.CreatedAt)
		if err != nil {
			ErrorLogger.Printf("Erreur lecture ligne commentaire: %v", err)
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// Ajouter ou supprimer des favoris
func ToogleFavorite(userID, recipeID int) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM Favorites WHERE user_id = ? AND recipe_id = ?", userID, recipeID).Scan(&exists)
	if err != nil {
		ErrorLogger.Printf("Erreur vérification favorite: %v", err)
		return err
	}
	if exists > 0 {
		_, err = db.Exec("DELETE FROM Favorites WHERE user_id = ? AND recipe_id = ?", userID, recipeID)
	} else {
		_, err = db.Exec("INSERT INTO Favorites (user_id, recipe_id) VALUES (?, ?)", userID, recipeID)
	}
	if err != nil {
		ErrorLogger.Printf("Erreur toggle favorite (user %d, recipe %d): %v", userID, recipeID, err)
	}
	return err
}

// Supprimer une recette + verification que l'utilisateur en est bien l'auteur
func DeleteRecipe(recipeID, userID int) error {
	res, err := db.Exec("DELETE FROM Recipes WHERE id = ? AND user_id = ?", recipeID, userID)
	if err != nil {
		ErrorLogger.Printf("Erreur suppression recette %d: %v", recipeID, err)
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Mettre à jour le pseudo ou la photo de profil
func UpdateUserProfile(userID int, username, photoPath string) error {
	_, err := db.Exec("UPDATE Users SET username = ?, photo_profil = ? WHERE id = ?", username, photoPath, userID)
	if err != nil {
		ErrorLogger.Printf("Erreur mise à jour profil (user %d): %v", userID, err)
	}
	return err
}
