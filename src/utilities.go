package yboost

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GenerateToken(pseudo string) string {
	expiration := time.Now().Add(1 * time.Hour)
	data := fmt.Sprintf("%s|%s", pseudo, expiration)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func GetUserFromSession(w http.ResponseWriter, r *http.Request) *UserInfos {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		return &UserInfos{Username: "Invité", IsConnected: false}
	}

	var user UserInfos
	query := `
        SELECT Users.username, Users.email 
        FROM Users 
        INNER JOIN Session ON Users.id = Session.user_id 
        WHERE Session.token = ?`

	err = db.QueryRow(query, cookie.Value).Scan(&user.Username, &user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Erreur vérification session: %v", err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		return &UserInfos{Username: "Invité", IsConnected: false}
	}

	user.IsConnected = true
	return &user
}

func SetFlash(w http.ResponseWriter, kind, message string) {
	value := url.QueryEscape(kind + "|" + message)
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    value,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetFlash(w http.ResponseWriter, r *http.Request) (string, string) {
	cookie, err := r.Cookie("flash")
	if err != nil {
		return "", ""
	}
	value, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		ClearFlash(w)
		return "", ""
	}
	parts := strings.SplitN(value, "|", 2)
	if len(parts) != 2 {
		ClearFlash(w)
		return "", ""
	}
	ClearFlash(w)
	return parts[0], parts[1]
}

func ClearFlash(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func LoadFlash(w http.ResponseWriter, r *http.Request, userInfos *UserInfos) {
	kind, message := GetFlash(w, r)
	userInfos.FlashType = kind
	userInfos.FlashMessage = message
	userInfos.PrintError = ""
}
