package models

import (
	u "../utils"
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title     string    `json:"title"`
	ImagePath string    `json:"image_path"`
	UserID    uint      `json:"user_id"`
	User      Account   `gorm:"foreignkey:UserID"`
	SectionID uint      `json:"section_id"`
	Section   Section   `gorm:"foreignkey:SectionID"`
	Comments  []Comment `gorm:"foreignkey:PostID"`
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
	GetDB().Preload("Section").Preload("User").Order("created_at").Find(&posts)
	return posts
}

func GetPost(u uint) *Post {

	post := &Post{}
	GetDB().Table("posts").Where("id = ?", u).First(post)
	if post.Title == "" { //User not found!
		return nil
	}
	return post
}

func GetPostWithComments(u uint) *Post {

	post := &Post{}
	GetDB().Preload("Comments").Preload("Comments.User").Preload("Section").Preload("User").Find(post, u)
	if post.Title == "" { //User not found!
		return nil
	}
	return post
}

func DeletePost(id int) {
	post := &Post{}
	db.First(&post, id)
	GetDB().Delete(&post)
}
