package controllers

import (
	"encoding/json"
	"fmt"
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
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	id, _ := strconv.Atoi(r.Header.Get("user_id"))
	file, fileheader, err := r.FormFile("file")
	if err != nil {
		fmt.Println("error-----" + err.Error())
	}
	requestParams := &models.RequestUserUpdate{
		FirstName:     r.FormValue("first_name"),
		LastName:      r.FormValue("last_name"),
		UserName:      r.FormValue("username"),
		ID:            id,
		File:          file,
		FileHeader:    fileheader,
		ProfileStatus: r.FormValue("profile_status"),
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.UpdateProfile(requestParams)
		u.Respond(w, resp)
		return
	}
}

var CheckUserName = func(w http.ResponseWriter, r *http.Request) {

	requestParams := &models.UserName{}
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
		resp := models.CheckUserName(requestParams)
		u.Respond(w, resp)
		return
	}
}
