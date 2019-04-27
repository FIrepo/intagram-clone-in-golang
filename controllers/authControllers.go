package controllers

import (
	"encoding/json"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
	"strconv"

	"github.com/chr4/pwgen"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	requestParams := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(requestParams) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(0, "Invalid request"))
		return
	}
	requestParams.VerificationToken = pwgen.Num(6)
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.CreateAccount(requestParams)
		u.Respond(w, resp)
		return
	}
}

var CheckMobileOTP = func(w http.ResponseWriter, r *http.Request) {
	requestParams := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(requestParams) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(0, "Invalid request"))
		return
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.CheckMobileOTP(requestParams)
		u.Respond(w, resp)
		return
	}
}

var UpdateProfile = func(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header.Get("user_id"))
	requestParams := &models.RequestUserUpdate{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		UserName:  r.FormValue("username"),
		ID:        id,
	}
	resp := models.UpdateProfile(requestParams)
	u.Respond(w, resp)
}

// // var Authenticate = func(w http.ResponseWriter, r *http.Request) {

// // 	account := &models.Account{}
// // 	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
// // 	if err != nil {
// 		u.Respond(w, u.Message(false, "Invalid request"))
// 		return
// 	}

// 	resp := models.Login(account.Email, account.Password)
// 	u.Respond(w, resp)
// }
