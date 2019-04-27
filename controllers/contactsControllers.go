package controllers

import (
	"encoding/json"
	"fmt"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
)

var CreateContact = func(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	formdata := r.MultipartForm // ok, no problem so far, read the Form data
	fmt.Println(formdata)
	_, ok := formdata.File["file"]
	if !ok {
		u.Respond(w, u.Message(0, "video file is not found"))
		return
	} else {
		files := formdata.File["file"][0] // grab the filenames
		fmt.Println(files)
		_, ok := formdata.File["thumb_file"]
		if !ok {
			u.Respond(w, u.Message(0, "Thumb file is not found"))
			return
		} else {
			thumbfiles := formdata.File["thumb_file"][0] // grab the filenames

			fmt.Println(files, thumbfiles)
			requestParams := models.RequestFile{
				File:      files,
				ThumbFile: thumbfiles,
				Language:  r.FormValue("language"),
				Category:  r.FormValue("category"),
				Title:     r.FormValue("title"),
			}
			objvalidate, haserror := u.ValidateObject(requestParams)
			if !haserror {
				u.Respond(w, objvalidate)
				return
			} else {
				resp := models.UploadFile(&requestParams)
				u.Respond(w, resp)
				return
			}
		}
	}
}

var VideoList = func(w http.ResponseWriter, r *http.Request) {
	requestParams := models.RequestList{}
	err := json.NewDecoder(r.Body).Decode(&requestParams) //decode the request body into struct and failed if any error occur
	if err != nil {
		fmt.Println(err.Error())
		u.Respond(w, u.Message(0, "Invalid request"))
		return
	}
	resp := models.VideoList(&requestParams)
	u.Respond(w, resp)
}

var ViewCount = func(w http.ResponseWriter, r *http.Request) {
	requestParams := models.RequestViewCountDelete{}
	err := json.NewDecoder(r.Body).Decode(&requestParams) //decode the request body into struct and failed if any error occur
	if err != nil {
		fmt.Println(err.Error())
		u.Respond(w, u.Message(0, "Invalid request"))
		return
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.ViewCount(&requestParams)
		u.Respond(w, resp)
	}
}
var DeleteVideo = func(w http.ResponseWriter, r *http.Request) {
	requestParams := models.RequestViewCountDelete{}
	err := json.NewDecoder(r.Body).Decode(&requestParams) //decode the request body into struct and failed if any error occur
	if err != nil {
		fmt.Println(err.Error())
		u.Respond(w, u.Message(0, "Invalid request"))
		return
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.DeleteVideo(&requestParams)
		u.Respond(w, resp)
	}
}

// var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {

// 	id := r.Context().Value("user").(uint)
// 	data := models.GetContacts(id)
// 	resp := u.Message(true, "success")
// 	resp["data"] = data
// 	u.Respond(w, resp)
// }
