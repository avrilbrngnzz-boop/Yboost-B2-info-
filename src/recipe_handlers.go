package yboost

import (
	"html/template"
	"net/http"
)
	
func createRecipeHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	user := GetUserFromSession(w, r)
	if !user.IsConnected {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}
	var userID int
	err := db.QueryRow("SELECT id FROM Users WHERE email = ?", user.Email).Scan(&userID)
	if err != nil {
		ErrorLogger.Printf("Erreur récupération ID user pour la création de la recette: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	ingredients := r.FormValue("ingredients")
	country := r.FormValue("country_of_origin")
	dishType := r.FormValue("dish_type")
	dietType := r.FormValue("diet_type")
	instructions := r.FormValue("instructions")
	photoURL := r.FormValue("photo_url")

	if title != "" && ingredients != "" && instructions != "" {
		err = CreateRecipe(userID, title, description, ingredients, country, dishType, dietType, instructions, photoURL)
		if err != nil {
			ErrorLogger.Printf("Erreur insertion recette: %v", err)
		}
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}



type PageData struct {
	User    *UserInfos
	Recipes []Recipe
}

func recipeHomeHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	user := GetUserFromSession(w, r)
	
	// Préparer les infos de l'utilisateur connecté
	currentUser := &UserInfos{
		Username:    user.Username,
		Email:       user.Email,
		IsConnected: user.IsConnected,
	}

	if user.IsConnected {
		pp, err := getUserPP(user.Email)
		if err == nil {
			currentUser.LoadedPP = pp
		}
	}

	// Récupérer les recettes depuis la BDD
	recipes, err := GetAllRecipes()
	if err != nil {
		ErrorLogger.Printf("Erreur chargement recettes: %v", err)
	}

	// Regrouper les données pour le template
	data := PageData{
		User:    currentUser,
		Recipes: recipes,
	}

	tmpl, err := template.ParseFiles("pages/home.html", "tmpl/header.html", "tmpl/footer.html")
	if err != nil {
		ErrorLogger.Printf("Erreur template home: %v", err)
		http.Error(w, "Erreur chargement page", http.StatusInternalServerError)
		return
	}
	
	tmpl.ExecuteTemplate(w, "home.html", data)
}