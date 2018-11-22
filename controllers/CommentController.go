package controllers

import (
	"net/http"
	"strconv"
	"time"

	"../models"
	u "../utils"
	"github.com/gorilla/mux"
)

var CreateComment = func(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	imagePath := ""
	if file != nil {
		imagePath = "/images/" + strconv.FormatInt(time.Now().UnixNano()/1000000, 10) + handler.Filename
		resp, isSuccess := u.UploadFile(file, handler, "./public"+imagePath)

		if isSuccess == false {
			u.RespondWithStatus(w, resp, http.StatusBadRequest)
			return
		}
		defer file.Close()
	}

	uid := r.Context().Value("user").(uint)
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		u.RespondWithStatus(w, u.Message(false, "format post id tidak sesuai"), http.StatusBadRequest)
		return
	}

	comment := &models.Comment{}
	comment.ImagePath = imagePath
	comment.Text = r.FormValue("text")
	comment.UserID = uid
	comment.PostID = postID
	resp := comment.Create()
	var httpStatus int
	if resp["status"] == false {
		httpStatus = http.StatusBadRequest
	} else {
		httpStatus = http.StatusCreated
	}

	u.RespondWithStatus(w, resp, httpStatus)
	return

}

var DeleteComment = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := r.Context().Value("user").(uint)
	commentID, _ := strconv.Atoi(vars["id"])
	comment := models.GetComment(uint(commentID))
	if comment == nil {
		u.RespondWithStatus(w, u.Message(false, "Komentar tidak ditemukan"), http.StatusNotFound)
		return
	}

	if comment.UserID != uid {
		u.RespondWithStatus(w, u.Message(false, "You are not authorized"), http.StatusUnauthorized)
		return
	}
	models.DeleteComment(commentID)
	resp := u.Message(true, "data berhasil dihapus")
	u.Respond(w, resp)
}

var UpVotesComment = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := r.Context().Value("user").(uint)
	commentID, _ := strconv.Atoi(vars["id"])
	comment := models.GetComment(uint(commentID))
	if comment == nil {
		u.RespondWithStatus(w, u.Message(false, "Komentar tidak ditemukan"), http.StatusNotFound)
		return
	}
	commentVotes := models.IsVotedComment(uid, uint(commentID))
	if commentVotes == nil {
		commentVotes = &models.CommentVotes{}
		commentVotes.UserID = uid
		commentVotes.CommentID = uint(commentID)
		commentVotes.Score = 1
		commentVotes.Create()
	} else if commentVotes.Score == 1 {
		u.RespondWithStatus(w, u.Message(false, "Anda sudah mengupvotes comment ini "), http.StatusBadRequest)
		return
	} else {
		models.UpdateCommentScore(commentVotes, 1)
	}
	u.Respond(w, u.Message(true, "upvote berhasil"))
}

var DownVotesComment = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := r.Context().Value("user").(uint)
	commentID, _ := strconv.Atoi(vars["id"])
	comment := models.GetComment(uint(commentID))
	if comment == nil {
		u.RespondWithStatus(w, u.Message(false, "Komentar tidak ditemukan"), http.StatusNotFound)
		return
	}
	commenttVotes := models.IsVotedComment(uid, uint(commentID))
	if commenttVotes == nil {
		commenttVotes = &models.CommentVotes{}
		commenttVotes.UserID = uid
		commenttVotes.CommentID = uint(commentID)
		commenttVotes.Score = -1
		commenttVotes.Create()
	} else if commenttVotes.Score == -1 {
		u.RespondWithStatus(w, u.Message(false, "Anda sudah mendownvote commentt ini "), http.StatusBadRequest)
		return
	} else {
		models.UpdateCommentScore(commenttVotes, -1)
	}
	u.Respond(w, u.Message(true, "downvote berhasil"))
}
