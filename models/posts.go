package models

import (
	"log"

	u "../utils"
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title     string      `json:"title"`
	ImagePath string      `json:"image_path"`
	UserID    uint        `json:"user_id"`
	User      Account     `gorm:"foreignkey:UserID"`
	SectionID uint        `json:"section_id"`
	Section   Section     `gorm:"foreignkey:SectionID"`
	Comments  []Comment   `gorm:"foreignkey:PostID"`
	PostVotes []PostVotes `gorm:"foreignkey:PostID"`
	Score     int         `gorm:"-"`
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
	//GetDB().Preload("Section").Preload("User").Order("created_at").Find(&posts)
	rows, _ := GetDB().Raw("SELECT  posts.id,posts.title,posts.image_path,  COALESCE(SUM(post_votes.score),0) ,posts.created_at,posts.deleted_at,sections.name ,sections.id AS 'section.id',	accounts.email ,accounts.id AS 'user.id' FROM posts LEFT JOIN post_votes ON posts.id = post_votes.post_id  LEFT JOIN sections ON sections.id = posts.section_id LEFT JOIN accounts ON accounts.id = posts.user_id  WHERE posts.deleted_at is NULL GROUP BY posts.title").Rows()
	defer rows.Close()
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID,
			&post.Title,
			&post.ImagePath,
			&post.Score,
			&post.CreatedAt,
			&post.DeletedAt,
			&post.Section.Name,
			&post.Section.ID,
			&post.User.Email,
			&post.User.ID)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, &post)
	}
	return posts
}

func GetPost(u uint) *Post {

	post := &Post{}
	GetDB().Where("id = ?", u).First(post)
	if post.Title == "" { //User not found!
		return nil
	}
	return post
}

func GetPostWithComments(u uint) *Post {

	post := Post{}
	GetDB().Preload("Section").Preload("User").Find(&post, u)
	row := GetDB().Raw("SELECT SUM(score) FROM post_votes WHERE post_id = ? ", u).Row()
	row.Scan(&post.Score)
	var comments []Comment
	//GetDB().Preload("Section").Preload("User").Order("created_at").Find(&posts)
	rows, _ := GetDB().Raw("SELECT  comments.id,comments.text,comments.image_path,  COALESCE(SUM(comment_votes.score),0) ,comments.created_at,comments.deleted_at,accounts.email ,accounts.id AS 'user.id' FROM comments LEFT JOIN comment_votes ON comments.id = comment_votes.comment_id   LEFT JOIN accounts ON accounts.id = comments.user_id  WHERE comments.deleted_at is NULL GROUP BY comments.id").Rows()
	defer rows.Close()
	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(&comment.ID,
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
	if post.Title == "" {
		return nil
	}
	post.Comments = comments
	return &post
}

func DeletePost(id int) {
	post := &Post{}
	db.First(&post, id)
	GetDB().Delete(&post)
}
