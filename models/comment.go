package models

import (
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

func DeleteComment(id int) {
	comment := &Comment{}
	GetDB().Delete(&comment, id)
}
