package yboost

import (
	"database/sql"
	"time"
)

type Post struct {
	ID			int
	UserID		int
	Author		string
	AuthorPP	string
	Title		string
	Content		string
	CreatedAt	time.Time
	LikesCount 	int
}

type Comment struct {
	ID			int
	PostID		int
	UserID		int
	Author		string
	AuthorPP	string
	Content		string
	CreatedAt	time.Time
}

func CreatePost(userID int, title, content string) error {
	_, err := db.Exec("INSERT INTO Posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		ErrorLogger.Printf("Erreur création du post (user %d): %v", userID, err)
	}
	return err
}

func CreateComment(userID, postID int, content string) error {
	_, err := db.Exec("INSERT INTO Comments (user_id, post_id, content) VALUES (?, ?, ?)", userID, postID, content)
	if err != nil {
		ErrorLogger.Printf("Erreur création du commentaire (user %d, post %d): %v", userID, postID, err)
	}
	return err
}

//Ajout like s'il n'existe pas et sinon le supprimer
func ToogleLike(userID, postID int) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM Likes Where user_id = ? AND post_id = ?", userID, postID).Scan(&exists)
	if err != nil {
		ErrorLogger.Printf("Erreur vérification like %v", err)
		return err
	}
	if exists > 0 {
		_, err = db.Exec("DELETE FROM Likes WHERE user_id = ? AND post_id = ?", userID, postID)
	} else {
		_, err = db.Exec("INSERT INTO Likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	}
	if err != nil {
		ErrorLogger.Printf("Erreur toggle like (user %d, post %d): %v", userID, postID, err)
	}
	return err
}

//Verification que le post appartient bien à l'utilisateur avant suppression
func DeletePost(postID, userID int) error {
	res, err := db.Exec("DELETE FROM Posts WHERE id = ? AND user_id = ?", postID, userID)
	if err != nil {
		ErrorLogger.Printf("Erreur suppression post %d: %v", postID, err)
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows //Pas de post trouvé pour cet utilisateur
	}
	return nil
}

//mettre a jour le pseudo ou la photo de profil
func UpdateUserProfile(userID int, username, photoPath string) error {
	_, err := db.Exec("UPDATE Users SET username = ?, photo_profil = ? WHERE id = ?", username, photoPath, userID)
	if err != nil {
		ErrorLogger.Printf("Erreur mise à jour profil (user %d): %v", userID, err)
	}
	return err
}