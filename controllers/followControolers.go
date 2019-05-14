package controllers

import (
	"fmt"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
)

func FollowingRequest(w http.ResponseWriter, r *http.Request) {
	requestParams := &models.RequestFollowing{
		ID:              r.Header.Get("user_id"),
		FollowingUserID: r.FormValue("following_user_id"),
		Status:          r.FormValue("status"),
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.FollowingRequest(requestParams)
		u.Respond(w, resp)
		return
	}
}

func FollowGet(w http.ResponseWriter, r *http.Request) {
	requestParams := &models.RequestFollowing{
		ID:              r.Header.Get("user_id"),
		Status:          r.FormValue("status"),
		FollowingUserID: "1",
	}
	if requestParams.Status != "" {
		if requestParams.Status != "0" && requestParams.Status != "1" {
			u.Respond(w, u.Message(0, "Invalid request."))
			return
		}
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.FollowGet(requestParams)
		u.Respond(w, resp)
		return
	}
}

func FollowersAction(w http.ResponseWriter, r *http.Request) {
	requestParams := &models.RequestFollowerAction{
		UserID:         r.Header.Get("user_id"),
		FollowerID:     r.FormValue("follower_id"),
		RejectAcceptID: r.FormValue("reject_accept_id"),
		RequestUserID:  r.FormValue("requested_user_id"),
	}
	fmt.Println(requestParams.RejectAcceptID)
	if requestParams.RejectAcceptID != "1" && requestParams.RejectAcceptID != "2" {
		u.Respond(w, u.Message(0, "Invalid request."))
		return
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.FollowersAction(requestParams)
		u.Respond(w, resp)
		return
	}
}

func GetUserList(w http.ResponseWriter, r *http.Request) {
	requestParams := &models.RequestUserList{
		UserID: r.Header.Get("user_id"),
		Name:   r.FormValue("name"),
	}
	objvalidate, haserror := u.ValidateObject(requestParams)
	if !haserror {
		u.Respond(w, objvalidate)
		return
	} else {
		resp := models.GetUserList(requestParams)
		u.Respond(w, resp)
		return
	}
}
