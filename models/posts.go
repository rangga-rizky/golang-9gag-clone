package models

import (
	u "../utils"
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title     string `json:"title"`
	ImagePath string `json:"image_path"`
	UserID    uint   `json:"user_id"`
	SectionID uint   `json:"section_id"`
}

func (post *Post) Validate() (map[string]interface{}, bool) {

	if post.Title == "" {
		return u.Message(false, "Judul Harus ada"), false
	}

	if post.SectionID <= 0 {
		return u.Message(false, "Section ID Harus ada"), false
	}

	return u.Message(true, "success"), true
}

func (post *Post) Create() map[string]interface{} {

	if resp, ok := post.Validate(); !ok {
		return resp
	}

	GetDB().Create(post)

	resp := u.Message(true, "success")
	resp["post"] = post
	return resp
}

func GetPosts() []*Post {

	var posts []*Post
	GetDB().Table("posts").Find(&posts)
	return posts
}
