package controllers

import (
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
)

func UploadPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		u.Respond(w, u.Message(0, "File is not found"))
		return
	}
	// file, fileheader, err := r.FormFile("file")
	// if err != nil {
	// 	u.Respond(w, u.Message(0, "File is not found"))
	// 	return
	// }
	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["multiplefiles"] // grab the filenames

	requestParams := &models.RequestUploadPost{
		Caption:      r.FormValue("caption"),
		Latitude:     r.FormValue("latitude"),
		Longitude:    r.FormValue("longitude"),
		ID:           r.Header.Get("user_id"),
		FileHeader:   files,
		LocalID:      r.Header.Get("user_local_id"),
		LocationName: r.FormValue("location_name"),
	}
	if len(requestParams.FileHeader) == 0 {
		u.Respond(w, u.Message(0, "File is not found"))
		return
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.UploadPost(requestParams)
		u.Respond(w, resp)
		return
	}
}

func GetListOfPost(w http.ResponseWriter, r *http.Request) {
	requestParams := &models.RequestListOfPost{
		ID: r.Header.Get("user_id"),
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.GetListOfPost(requestParams)
		u.Respond(w, resp)
		return
	}
}
