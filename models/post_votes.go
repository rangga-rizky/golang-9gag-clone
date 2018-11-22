package models

import (
	u "../utils"
	"github.com/jinzhu/gorm"
)

type PostVotes struct {
	gorm.Model
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
	Score  int  `json:"score"`
}

func (postVotes *PostVotes) Validate() (map[string]interface{}, bool) {

	if postVotes.PostID <= 0 {
		return u.Message(false, "Post ID Harus ada"), false
	}

	return u.Message(true, "success"), true
}

func (postVotes *PostVotes) Create() map[string]interface{} {

	if resp, ok := postVotes.Validate(); !ok {
		return resp
	}

	GetDB().Create(postVotes)

	resp := u.Message(true, "success")
	resp["postVotes"] = postVotes
	return resp
}

func IsVotedPost(uid uint, postID uint) *PostVotes {

	postVotes := &PostVotes{}
	GetDB().Where("user_id = ?", uid).Where("post_id = ?", postID).First(postVotes)
	if postVotes.ID == 0 {
		return nil
	}
	return postVotes
}

func UpdatePostScore(postVotes *PostVotes, score int) {
	postVotes.Score = score
	GetDB().Save(&postVotes)
}
