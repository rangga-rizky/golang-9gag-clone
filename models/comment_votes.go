package models

import (
	u "../utils"
	"github.com/jinzhu/gorm"
)

type CommentVotes struct {
	gorm.Model
	UserID    uint `json:"user_id"`
	CommentID uint `json:"comment_id"`
	Score     int  `json:"score"`
}

func (commentVotes *CommentVotes) Validate() (map[string]interface{}, bool) {

	if commentVotes.CommentID <= 0 {
		return u.Message(false, "Comment ID Harus ada"), false
	}

	return u.Message(true, "success"), true
}

func (commentVotes *CommentVotes) Create() map[string]interface{} {

	if resp, ok := commentVotes.Validate(); !ok {
		return resp
	}

	GetDB().Create(commentVotes)

	resp := u.Message(true, "success")
	resp["commentVotes"] = commentVotes
	return resp
}

func IsVotedComment(uid uint, commentID uint) *CommentVotes {

	commentVotes := &CommentVotes{}
	GetDB().Where("user_id = ?", uid).Where("comment_id = ?", commentID).First(commentVotes)
	if commentVotes.ID == 0 {
		return nil
	}
	return commentVotes
}

func UpdateCommentScore(commentVotes *CommentVotes, score int) {
	commentVotes.Score = score
	GetDB().Save(&commentVotes)
}
