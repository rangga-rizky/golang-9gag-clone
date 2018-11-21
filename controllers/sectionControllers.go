package controllers

import (
	"net/http"

	"../models"
	u "../utils"
)

var GetSections = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetSections()
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)

}
