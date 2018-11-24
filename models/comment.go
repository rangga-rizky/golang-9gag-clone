package models

import (
	"log"

	u "../utils"
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	User      Account `gorm:"foreignkey:UserID"`
	PostID    int     `json:"post_id"`
	ImagePath string  `json:"image_path"`
	Text      string  `json:"text"`
	Score     int     `gorm:"-"`
}

func (comment *Comment) Validate() (map[string]interface{}, bool) {

	if comment.Text == "" && comment.ImagePath == "" {
		return u.Message(false, "Harus ada salah satu text/gambar"), false
	}

	if comment.PostID <= 0 {
		return u.Message(false, "Post ID Harus ada"), false
	}

	return u.Message(true, "success"), true
}

func (comment *Comment) Create() map[string]interface{} {

	if resp, ok := comment.Validate(); !ok {
		return resp
	}
	GetDB().Create(comment)
	resp := u.Message(true, "success")
	resp["comment"] = comment
	return resp
}

func GetComment(u uint) *Comment {

	comment := &Comment{}
	GetDB().Where("id = ?", u).First(comment)
	if comment.ID == 0 {
		return nil
	}
	return comment
}

func getCommentsByPost(postID uint) []Comment {
	var comments []Comment
	rows, _ := GetDB().Raw("SELECT  comments.id,comments.post_id,comments.text,comments.image_path,  COALESCE(SUM(comment_votes.score),0) ,comments.created_at,comments.deleted_at,accounts.email ,accounts.id AS 'user.id' FROM comments LEFT JOIN comment_votes ON comments.id = comment_votes.comment_id   LEFT JOIN accounts ON accounts.id = comments.user_id  WHERE comments.deleted_at is NULL AND comments.post_id = ? GROUP BY comments.id", postID).Rows()
	defer rows.Close()
	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(&comment.ID,
			&comment.PostID,
			&comment.Text,
			&comment.ImagePath,
			&comment.Score,
			&comment.CreatedAt,
			&comment.DeletedAt,
			&comment.User.Email,
			&comment.User.ID)
		if err != nil {
			log.Fatal(err)
		}
		comments = append(comments, comment)
	}

	return comments
}

func DeleteComment(id int) {
	comment := &Comment{}
	GetDB().Delete(&comment, id)
}
