package models

import (
	"fmt"
	u "go-contacts/utils"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/rs/xid"
	"gopkg.in/guregu/null.v3/zero"
)

type (
	RequestUploadPost struct {
		// File       multipart.File          `json:"file" validate:"required"`
		FileHeader   []*multipart.FileHeader `json:"fileheader"`
		Caption      string                  `json:"caption"`
		ID           string                  `json:"id" validate:"required"`
		Latitude     string                  `json:"latitude"`
		Longitude    string                  `json:"longitude"`
		LocalID      string                  `json:"local_id" validate:"required"`
		LocationName string                  `json:"location_name" db:"location_name"`
	}
	RequestListOfPost struct {
		ID string `json:"id" validate:"required"`
	}
	SelectGetListOfPost struct {
		ID             int         `json:"id" db:"id"`
		UserID         int         `json:"user_id" db:"user_id"`
		Caption        zero.String `json:"caption" db:"caption"`
		Latitude       zero.String `json:"latitude" db:"latitude"`
		Longitude      zero.String `json:"longitude" db:"longitude"`
		PostPath       zero.String `json:"post_path" db:"post_path"`
		CreatedAt      int64       `json:"created_at" db:"created_at"`
		Name           zero.String `json:"name" db:"name"`
		ProfilePicture zero.String `json:"profile_picture" db:"profile_picture"`
		LocationName   zero.String `json:"location_name" db:"location_name"`
	}
)

func UploadPost(requestParams *RequestUploadPost) map[string]interface{} {
	resp := make(map[string]interface{})
	conn := PgCon()
	LastInsertID := 0
	err := conn.QueryRow(`INSERT INTO public.post(
		user_id, caption, latitude, longitude,location_name)
		VALUES ($1, $2, $3, $4,$5)  RETURNING id;`, requestParams.ID, requestParams.Caption, requestParams.Latitude, requestParams.Longitude, requestParams.LocationName).Scan(&LastInsertID)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	fmt.Println(LastInsertID)
	go MultipleFile(LastInsertID, requestParams.FileHeader, requestParams.LocalID)
	resp = u.Message(1, "successfully uploaded post.")
	conn.Close()
	return resp
}

func MultipleFile(lastID int, files []*multipart.FileHeader, localID string) {
	conn := PgCon()
	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error(), ":=== file number ", i, "   code line Multiple File")
			return
		}
		guid := xid.New()
		post := guid.String() + filepath.Ext(files[i].Filename)
		filePath := os.Getenv("user_profile") + localID + "/"
		_, err = conn.Exec(`INSERT INTO public.post_images(
							post_image, post_id)
							VALUES ($1,$2);`, filePath+post, int(lastID))
		if err != nil {
			conn.Close()
			fmt.Println(err.Error(), ":=== file number ", i, " query error  code line Multiple File")
			return
		}
		SaveFileUser(file, post, filePath)
	}
	conn.Close()
}

func GetListOfPost(requestParams *RequestListOfPost) map[string]interface{} {
	resp := make(map[string]interface{})
	conn := PgCon()
	conSelectGetListOfPost := make([]SelectGetListOfPost, 0)
	err := conn.Select(&conSelectGetListOfPost, `SELECT post.id as id, user_id, caption, latitude, longitude, post.created_at,location_name,post_images.post_image as post_path,concat(first_name,' ',last_name) as name,
	CASE 
		when profile_picture != '' then 
		concat('`+os.Getenv("Api_path")+":"+os.Getenv("PORT")+"/"+os.Getenv("user_profile")+`',local_id,'/',profile_picture)
		else 
		''
		end as profile_picture
		 FROM public.post 
		left join user_master on user_id = user_master.id
		left join (select array_agg(concat('`+os.Getenv("Api_path")+":"+os.Getenv("PORT")+"/"+`',post_image)) as post_image,post_id from post_images  group by post_id) as post_images on 
		post.id = post_images.post_id where user_id in (select following_id from public.following where user_id = $1 
    and reject_accept_key=1) or user_id = $1 order by created_at desc;`, requestParams.ID)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	resp = u.Message(1, "successfully uploaded post.")
	resp["data"] = conSelectGetListOfPost
	conn.Close()
	return resp
}
