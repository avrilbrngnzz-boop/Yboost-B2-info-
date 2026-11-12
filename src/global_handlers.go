package yboost

import (
	"html/template"
	"net/http"
)

func initHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	user := GetUserFromSession(w, r)
	userInfos.Username = user.Username
	userInfos.Email = user.Email
	LoadFlash(w, r, user)

	pp, err := getUserPP(user.Email)
	if err == nil {
		user.LoadedPP = pp
	} else {
		user.LoadedPP = ""
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func homeHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, _ := template.ParseFiles("pages/home.html", "tmpl/header.html", "tmpl/footer.html")

	tmpl.ExecuteTemplate(w, "home.html", nil)
}
