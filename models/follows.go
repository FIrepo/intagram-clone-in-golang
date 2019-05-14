package models

import (
	"fmt"
	u "go-contacts/utils"
	"os"

	"gopkg.in/guregu/null.v3/zero"
)

type (
	RequestFollowing struct {
		ID              string `json:"user_id" validate:"required"`
		FollowingUserID string `json:"request_user_id"  validate:"required"`
		Status          string `json:"status"`
	}
	RequestUserList struct {
		Name   string `json:"name" validate:"required"`
		UserID string `json:"user_id" validate:"required"`
	}
	RequestFollowerAction struct {
		UserID         string `json:"user_id" validate:"required"`
		FollowerID     string `json:"follower_id"  validate:"required"`
		RejectAcceptID string `json:"reject_accept_id"  validate:"required"`
		RequestUserID  string `json:"request_user_id" validate:"required"`
	}
	PrcSelectFollowerList struct {
		FollowerID      int         `json:"follower_id" db:"follower_id"`
		RequestedUserID int         `json:"requested_user_id" db:"requested_user_id"`
		FirstName       zero.String `json:"first_name" db:"first_name"`
		LastName        zero.String `json:"last_name" db:"last_name"`
		ProfilePicture  zero.String `json:"profile_picture" db:"profile_picture"`
	}
	PrcSelectUserList struct {
		UserID         int         `json:"user_id" db:"user_id"`
		FirstName      zero.String `json:"first_name" db:"first_name"`
		LastName       zero.String `json:"last_name" db:"last_name"`
		ProfilePicture zero.String `json:"profile_picture" db:"profile_picture"`
	}
)

func FollowingRequest(requestParams *RequestFollowing) map[string]interface{} {
	resp := make(map[string]interface{})
	conn := PgCon()
	_, err := conn.Exec(`INSERT INTO following(
		following_id, user_id)
		VALUES ($1, $2);`, requestParams.FollowingUserID, requestParams.ID)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	_, err = conn.Exec(`INSERT INTO public.followers(
			follower_id,user_id)
			VALUES ($2, $1);`, requestParams.FollowingUserID, requestParams.ID)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	resp = u.Message(1, "successfully request sent.")
	conn.Close()
	return resp
}

// if 1 = accept 2 = delete --------------
func FollowersAction(requestParams *RequestFollowerAction) map[string]interface{} {
	resp := make(map[string]interface{})
	conn := PgCon()
	if requestParams.RejectAcceptID == "2" {
		_, err := conn.Exec(`delete from followers where id = $1;`, requestParams.FollowerID)
		if err != nil {
			conn.Close()
			return u.Message(0, err.Error())
		}
		_, err = conn.Exec(`delete from following where following_id = $1 and user_id = $2`, requestParams.UserID, requestParams.RequestUserID)
		if err != nil {
			conn.Close()
			return u.Message(0, err.Error())
		}
	} else if requestParams.RejectAcceptID == "1" {
		_, err := conn.Exec(`update followers set reject_accept_key = 1 where id = $1;`, requestParams.FollowerID)
		if err != nil {
			conn.Close()
			return u.Message(0, err.Error())
		}
		_, err = conn.Exec(`update following set reject_accept_key = 1 where following_id = $1 and user_id = $2`, requestParams.UserID, requestParams.RequestUserID)
		if err != nil {
			conn.Close()
			return u.Message(0, err.Error())
		}
	}
	resp = u.Message(1, "successfully uploaded post.")
	conn.Close()
	return resp
}

func FollowGet(requestParams *RequestFollowing) map[string]interface{} {
	resp := make(map[string]interface{})
	conn := PgCon()
	SelectFollowerList := make([]PrcSelectFollowerList, 0)
	fmt.Println(requestParams)
	if requestParams.Status == "" {
		requestParams.Status = "0"
	}
	err := conn.Select(&SelectFollowerList, `select followers.id as follower_id,follower_id as requested_user_id,first_name,last_name,
		CASE 
		when profile_picture != '' then 
		concat('`+os.Getenv("Api_path")+":"+os.Getenv("PORT")+"/"+os.Getenv("user_profile")+`',local_id,'/',profile_picture)
		else 
		''
		end as profile_picture
		from followers left join user_master on user_master.id = follower_id where reject_accept_key = $2 and user_id =$1;`, requestParams.ID, requestParams.Status)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	resp = u.Message(1, "successfully uploaded post.")
	resp["data"] = SelectFollowerList
	conn.Close()
	return resp
}

func GetUserList(requestParams *RequestUserList) map[string]interface{} {
	resp := make(map[string]interface{})
	conn := PgCon()
	SelectUserList := make([]PrcSelectUserList, 0)
	err := conn.Select(&SelectUserList, `select id as user_id,first_name,last_name,
		CASE 
		when profile_picture != '' then 
		concat('`+os.Getenv("Api_path")+":"+os.Getenv("PORT")+"/"+os.Getenv("user_profile")+`',local_id,'/',profile_picture)
		else 
		''
		end as profile_picture
		from user_master where concat(first_name,' ',last_name) !='' and username !='' and id != $1 and (LOWER(concat(first_name,' ',last_name)) like '`+requestParams.Name+`%' or username like '`+requestParams.Name+`%');`, requestParams.UserID)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	resp = u.Message(1, "successfully uploaded post.")
	resp["data"] = SelectUserList
	conn.Close()
	return resp
}
