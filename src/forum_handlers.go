package yboost

import (
	"html/template"
	"net/http"
)
	
func createPostHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
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
		ErrorLogger.Printf("Erreur récupération ID user pour la création du post: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	if title != "" && content != "" {
		err = CreatePost(userID, title, content)
		if err != nil {
			ErrorLogger.Printf("Erreur insertion post: %v, err")
		}
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func forumHomeHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	user := GetUserFromSession(w, r)
	userInfos.Username = user.Username
	userInfos.Email = user.Email
	userInfos.IsConnected = user.IsConnected

	if user.IsConnected {
		pp, err := getUserPP(user.Email)
		if err == nil {
			userInfos.LoadedPP = pp
		}
	}
	posts, err := GetAllPosts()
	if err != nil {
		ErrorLogger.Printf("Erreur chargement posts: %v", err)
	}
	tmpl, err := template.ParseFiles("pages/home.html","tmpl/header.html", "tmpl/footer.html")
	if err != nil {
		ErrorLogger.Printf("Erreur template home: %v", err)
		http.Error(w, "Erreur chargement page", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "home.html", posts)
}