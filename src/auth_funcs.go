package yboost

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func checkPasswordCharacters(userInfos *UserInfos) {
	var allowedCharacters = regexp.MustCompile(`^[\x21-\x7E]+$`)
	var CPC_hasUpper = regexp.MustCompile(`[A-Z]`)
	var CPC_hasLower = regexp.MustCompile(`[a-z]`)
	var CPC_hasDigit = regexp.MustCompile(`[0-9]`)
	var hasSpecial = regexp.MustCompile(`[!"#$%&'()*+,\-./:;<=>?@[\\\]^_{|}~]`)

	if !allowedCharacters.MatchString(userInfos.Password) ||
		!CPC_hasUpper.MatchString(userInfos.Password) ||
		!CPC_hasLower.MatchString(userInfos.Password) ||
		!CPC_hasDigit.MatchString(userInfos.Password) ||
		!hasSpecial.MatchString(userInfos.Password) {
		userInfos.PrintError = "La composition du mot de passe ne respecte pas les critères attendus. Veuillez réessayer."
	}
}

func dataLoginCheck(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	var comparepassword_hash string
	var userID string

	if strings.Contains(userInfos.Email_Username, "@") {
		userInfos.Email = userInfos.Email_Username
		db.QueryRow("SELECT id, username, password_hash FROM Users WHERE email=?", userInfos.Email).Scan(&userID, &userInfos.Username, &comparepassword_hash)
	} else {
		userInfos.Username = userInfos.Email_Username
		db.QueryRow("SELECT id, email, password_hash FROM Users WHERE username=?", userInfos.Username).Scan(&userID, &userInfos.Email, &comparepassword_hash)
	}
	userInfos.DBid = userID

	if bcrypt.CompareHashAndPassword([]byte(comparepassword_hash), []byte(userInfos.Password)) == nil {
		userInfos.PrintError = ""
		sessionToken := GenerateToken(userInfos.Email)

		db.Exec("DELETE FROM Session WHERE user_id = ?", userID)
		_, err := db.Exec("INSERT INTO Session (user_id, token) VALUES (?, ?)", userID, sessionToken)
		if err != nil {
			fmt.Println("Erreur création session DB:", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   86400,
		})

		userInfos.Password = ""
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		userInfos.PrintError = "Email ou mot de passe incorrect."
		userInfos.Password = ""
		SetFlash(w, "error", userInfos.PrintError)
		println(userInfos.PrintError)
		userInfos.PrintError = ""
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
	}
}

func dataRegisterSend(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	userInfos.PrintError = ""

	checkUsernameChar := "!\"#$%&'()*+,-./:;<=>?@[\\]^ ` {|}~€£¥©®™§"
	var usernameCharIsOk = true
	for _, charac := range userInfos.Username {
		if strings.ContainsRune(checkUsernameChar, charac) {
			usernameCharIsOk = false
			break
		}
	}

	if !usernameCharIsOk {
		userInfos.PrintError = "Le seul caractère spécial autorisé pour le nom d'utilisateur est _ . Veuillez réessayer."
		userInfos.Password = ""
		userInfos.ConfirmPassword = ""
		SetFlash(w, "error", userInfos.PrintError)
		println(userInfos.PrintError)
		userInfos.PrintError = ""
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	if !emailRegex.MatchString(userInfos.Email) {
		userInfos.Password = ""
		userInfos.ConfirmPassword = ""
		userInfos.PrintError = "Le format du mail est incorrecte. Veuillez réessayer."
		SetFlash(w, "error", userInfos.PrintError)
		println(userInfos.PrintError)
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	checkPasswordCharacters(userInfos)
	if userInfos.PrintError != "" {
		userInfos.Password = ""
		userInfos.ConfirmPassword = ""
		SetFlash(w, "error", userInfos.PrintError)
		println(userInfos.PrintError)
		userInfos.PrintError = ""
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	if len(userInfos.Password) < 12 {
		userInfos.PrintError = "La taille du mot de passe doit être d'au moins 12 caractères."
		userInfos.Password = ""
		userInfos.ConfirmPassword = ""
		SetFlash(w, "error", userInfos.PrintError)
		println(userInfos.PrintError)
		userInfos.PrintError = ""
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	if userInfos.Password != userInfos.ConfirmPassword {
		userInfos.PrintError = "Les mots de passe ne correspondent pas."
		userInfos.Password = ""
		userInfos.ConfirmPassword = ""
		SetFlash(w, "error", userInfos.PrintError)
		println(userInfos.PrintError)
		userInfos.PrintError = ""
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(userInfos.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	res, err := db.Exec("INSERT INTO Users (username, email, password_hash) VALUES (?, ?, ?)", userInfos.Username, userInfos.Email, string(password_hash))
	if err != nil {
		errMsg := err.Error()
		println(errMsg)
		if strings.Contains(errMsg, "email") && strings.Contains(errMsg, "unique") {
			userInfos.PrintError = "Cet email est déjà utilisé."
		} else if strings.Contains(errMsg, "username") && strings.Contains(errMsg, "unique") {
			userInfos.PrintError = "Ce nom d'utilisateur est déjà utilisé."
		} else {
			userInfos.PrintError = "Impossible de créer le compte. Veuillez réessayer plus tard."
		}
		userInfos.Password = ""
		userInfos.ConfirmPassword = ""
		SetFlash(w, "error", userInfos.PrintError)
		println(userInfos.PrintError)
		userInfos.PrintError = ""
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	lastID, _ := res.LastInsertId()
	token := GenerateToken(userInfos.Email)
	db.Exec("INSERT INTO Session (user_id, token) VALUES (?, ?)", lastID, token)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	db.QueryRow("SELECT id FROM Users WHERE email=?", userInfos.Email).Scan(&userInfos.DBid)

	userInfos.Password = ""
	userInfos.ConfirmPassword = ""
	userInfos.PrintError = ""
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
