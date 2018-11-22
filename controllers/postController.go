package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"../models"
	serializer "../serializers"
	u "../utils"
)

var CreatePost = func(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if file == nil {
		u.RespondWithStatus(w, u.Message(false, "File Gambar tidak boleh kosong"), http.StatusBadRequest)
		return
	}

	if err != nil {
		u.RespondWithStatus(w, u.Message(false, err.Error()), http.StatusBadRequest)
		return
	}
	defer file.Close()
	imagePath := "/images/" + strconv.FormatInt(time.Now().UnixNano()/1000000, 10) + handler.Filename
	resp, isSuccess := u.UploadFile(file, handler, "./public"+imagePath)

	if isSuccess == false {
		u.RespondWithStatus(w, resp, http.StatusBadRequest)
		return
	}
	uid := r.Context().Value("user").(uint)
	sectionID, err := strconv.Atoi(r.FormValue("section_id"))
	if err != nil {
		u.RespondWithStatus(w, u.Message(false, "format section id tidak sesuai"), http.StatusBadRequest)
		return
	}

	post := &models.Post{}
	post.ImagePath = imagePath
	post.Title = r.FormValue("title")
	post.UserID = uid
	post.SectionID = uint(sectionID)
	resp = post.Create()
	var httpStatus int
	if resp["status"] == false {
		httpStatus = http.StatusBadRequest
	} else {
		httpStatus = http.StatusCreated
	}

	u.RespondWithStatus(w, resp, httpStatus)
	return

}

var GetPosts = func(w http.ResponseWriter, r *http.Request) {
	posts := models.GetPosts()
	data, _ := serializer.CustomPostSerializer().TransformArray(posts)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var DeletePost = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := r.Context().Value("user").(uint)
	postID, _ := strconv.Atoi(vars["id"])
	post := models.GetPost(uint(postID))
	if post == nil {
		u.RespondWithStatus(w, u.Message(false, "Post tidak ditemukan"), http.StatusNotFound)
		return
	}

	if post.UserID != uid {
		u.RespondWithStatus(w, u.Message(false, "You are not authorized"), http.StatusUnauthorized)
		return
	}
	models.DeletePost(postID)
	resp := u.Message(true, "data berhasil dihapus")
	u.Respond(w, resp)
}

var GetPost = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["id"])
	post := models.GetPostWithComments(uint(postID))
	if post == nil {
		u.RespondWithStatus(w, u.Message(false, "Post tidak ditemukan"), http.StatusNotFound)
		return
	}
	data := serializer.CustomPostSerializer().WithComments().Transform(post)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var UpVotes = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := r.Context().Value("user").(uint)
	postID, _ := strconv.Atoi(vars["id"])
	postVotes := models.IsVoted(uid, uint(postID))
	if postVotes == nil {
		postVotes = &models.PostVotes{}
		postVotes.UserID = uid
		postVotes.PostID = uint(postID)
		postVotes.Score = 1
		postVotes.Create()
	} else if postVotes.Score == 1 {
		u.RespondWithStatus(w, u.Message(false, "Anda sudah mengupvotes post ini "), http.StatusBadRequest)
		return
	} else {
		models.UpdateScore(postVotes, 1)
	}
	u.Respond(w, u.Message(true, "upvote berhasil"))
}

var DownVotes = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := r.Context().Value("user").(uint)
	postID, _ := strconv.Atoi(vars["id"])
	postVotes := models.IsVoted(uid, uint(postID))
	if postVotes == nil {
		postVotes = &models.PostVotes{}
		postVotes.UserID = uid
		postVotes.PostID = uint(postID)
		postVotes.Score = -1
		postVotes.Create()
	} else if postVotes.Score == -1 {
		u.RespondWithStatus(w, u.Message(false, "Anda sudah mendownvote post ini "), http.StatusBadRequest)
		return
	} else {
		models.UpdateScore(postVotes, -1)
	}
	u.Respond(w, u.Message(true, "downvote berhasil"))
}
