package yboost

import (
	"html/template"
	"log"
	"net/http"
)

func connexionHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	tmpl, err := template.ParseFiles("pages/connexion.html", "tmpl/header.html", "tmpl/footer.html")
	if err != nil {
		log.Printf("Erreur template login : %v", err)
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "connexion.html", nil)
}

func checkloginHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	userInfos.Email_Username = r.FormValue("email_Username")
	userInfos.Password = r.FormValue("password")

	dataLoginCheck(w, r, userInfos)
}

func checkregisterHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	userInfos.Username = r.FormValue("username")
	userInfos.Email = r.FormValue("email")
	userInfos.Password = r.FormValue("password")
	userInfos.ConfirmPassword = r.FormValue("confirmpassword")
	dataRegisterSend(w, r, userInfos)
}

func editaccountHandler(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	tmpl, err := template.ParseFiles("pages/account.html", "tmpl/header.html", "tmpl/footer.html")
	if err != nil {
		log.Printf("Erreur template account : %v", err)
		return
	}

	tmpl.ExecuteTemplate(w, "account.html", nil)
}
