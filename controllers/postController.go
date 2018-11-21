package controllers

import (
	"net/http"
	"strconv"
	"time"

	"../models"
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
	imagePath := "/public/images/" + strconv.FormatInt(time.Now().UnixNano()/1000000, 10) + handler.Filename
	resp, isSuccess := u.UploadFile(file, handler, "."+imagePath)

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
	data := models.GetPosts()
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
