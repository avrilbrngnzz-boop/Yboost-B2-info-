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

func GetAllPosts() ([]Post, error) { //COALESCE permet d'eviter une erreur GO quand la photo de profil est NULL
	query := `
		SELECT Posts.id, Posts.user_id, Users.username, COALESCE(Users.photo_profil, ''), Posts.title, Posts.content, Posts.created_at, 
		       (SELECT COUNT(*) FROM Likes WHERE Likes.post_id = Posts.id) AS likes_count
		FROM Posts
		INNER JOIN Users ON Posts.user_id = Users.id
		ORDER BY Posts.created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		ErrorLogger.Printf("Erreur récupération posts: %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Author, &p.AuthorPP, &p.Title, &p.Content, &p.CreatedAt, &p.LikesCount)
		if err != nil {
			ErrorLogger.Printf("Erreur lecture ligne post: %v", err)
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetCommentsByPostID(postID int) ([]Comment, error) {
	query := `
		SELECT Comments.id, Comments.post_id, Comments.user_id, Users.username, COALESCE(Users.photo_profil, ''), Comments.content, Comments.created_at
		FROM Comments
		INNER JOIN Users ON Comments.user_id = Users.id
		WHERE Comments.post_id = ?
		ORDER BY Comments.created_at ASC
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		ErrorLogger.Printf("Erreur récupération commentaires (post %d): %v", postID, err)
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Author, &c.AuthorPP, &c.Content, &c.CreatedAt)
		if err != nil {
			ErrorLogger.Printf("Erreur lecture ligne commentaire: %v", err)
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}