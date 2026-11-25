package yboost

import (
	"log"
	"net/http"
)

const serverAddr = ":8080"

type UserInfos struct {
	Username        string
	Email           string
	Email_Username  string
	Password        string
	ConfirmPassword string
	LoadedPP        string
	PrintError      string
	FlashMessage    string
	FlashType       string
	DBid            string
	IsConnected     bool
}

var userInfos UserInfos

func StartServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		initHandler(w, r, &userInfos)
	})

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		// homeHandler(w, r, &userInfos)
		recipeHomeHandler(w, r, &userInfos)
	})

	http.HandleFunc("/createrecipe", func(w http.ResponseWriter, r *http.Request) {
		createRecipeHandler(w, r, &userInfos)
	})

	http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {
		connexionHandler(w, r, &userInfos)
	})
	
	http.HandleFunc("/checklogin", func(w http.ResponseWriter, r *http.Request) {
		checkloginHandler(w, r, &userInfos)
	})

	http.HandleFunc("/checkregister", func(w http.ResponseWriter, r *http.Request) {
		checkregisterHandler(w, r, &userInfos)
	})

	http.HandleFunc("/editaccount", func(w http.ResponseWriter, r *http.Request) {
		if userInfos.Email != "" {
			userInfos.LoadedPP, _ = getUserPP(userInfos.Email)
		}
		editaccountHandler(w, r, &userInfos)
	})

	log.Printf("Server listening on https://localhost%s", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
