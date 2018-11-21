package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"../models"
	u "../utils"
	"github.com/gorilla/mux"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	VerificationToken := u.RandStringBytes(15)
	account.VerificationToken = VerificationToken
	resp := account.Create() //Create account
	if resp["status"] == true {
		u.SendMail(account.Email, "Account Verification", VerificationToken)
	}
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondWithStatus(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	resp := models.Login(account.Email, account.Password)
	status := http.StatusOK
	if resp["status"] == false {
		status = http.StatusUnauthorized
	}

	u.RespondWithStatus(w, resp, status)

}

var VerifyEMail = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["uid"])
	user := models.GetUser(uint(id))

	if user.VerificationToken == vars["token"] {
		models.VerifyUser(uint(id))
		u.RespondWithStatus(w, u.Message(true, "Verifikasi Berhasil anda dapat login sekarang"), http.StatusOK)
	} else {
		u.RespondWithStatus(w, u.Message(false, "Token not Valid"), http.StatusNotFound)
	}

}
