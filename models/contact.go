package models

import (
	"fmt"
	u "go-contacts/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/xid"
	"gopkg.in/guregu/null.v3/zero"
)

type (
	RequestFile struct {
		File      *multipart.FileHeader `json:"file" validate:"required"`
		ThumbFile *multipart.FileHeader `json:"thumbfile" validate:"required"`
		Title     string                `json:"title" validate:"required"`
		Language  string                `json:"language" validate:"required"`
		Category  string                `json:"cetgory" validate:"required"`
	}
	RequestList struct {
		Category  string `json:"category"`
		Language  string `json:"language"`
		CreatedOn string `json:"created_at"`
		Title     string `json:"title"`
	}
	RequestViewCountDelete struct {
		ID int `json:"id" validate:"required"`
	}
	ResponseStruct struct {
		ID        int         `json:"id" db:"id"`
		Title     zero.String `json:"title" db:"title"`
		VideoUrl  zero.String `json:"video_url" db:"video_url"`
		ThumbUrl  zero.String `json:"thumb_url" db:"thumb_url"`
		Language  zero.String `json:"language" db:"language"`
		CreatedAt int64       `json:"created_at" db:"created_at"`
		ViewCount zero.Int    `json:"view_count" db:"view_count"`
		Category  zero.String `json:"category" db:"category"`
	}
)

func UploadFile(requestParams *RequestFile) map[string]interface{} {
	resp := make(map[string]interface{})
	VideoFile, err := requestParams.File.Open()
	if err != nil {
		return u.Message(0, err.Error())
	}
	ThumbFile, err := requestParams.ThumbFile.Open()
	if err != nil {
		return u.Message(0, err.Error())
	}
	conn := PgCon()
	var video, thumburl string
	// ApiPath := os.Getenv("Api_path") + ":" + os.Getenv("PORT") + "/" + os.Getenv("video_path")
	guid := xid.New()
	video = guid.String() + filepath.Ext(requestParams.File.Filename)
	thumburl = guid.String() + "_thumb" + filepath.Ext(requestParams.ThumbFile.Filename)
	stmt, err := conn.Exec("INSERT INTO public.videostatus(title, video_url, thumb_url, category, language) VALUES ($1,$2,$3,$4,$5)", requestParams.Title, os.Getenv("video_path")+video, os.Getenv("video_path")+thumburl, requestParams.Category, requestParams.Language)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	go SaveFile(VideoFile, video)
	go SaveFile(ThumbFile, thumburl)
	fmt.Println(stmt.RowsAffected())
	resp = u.Message(1, "successfully upload video.")
	conn.Close()
	return resp
}
func SaveFile(File multipart.File, FileName string) {
	f, err := os.OpenFile("./"+os.Getenv("video_path")+FileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("error :- ", err)
	}
	defer f.Close()
	io.Copy(f, File)
}
func VideoList(requestParams *RequestList) map[string]interface{} {
	resp := make(map[string]interface{})
	data := make(map[string]interface{})
	conn := PgCon()
	resp = u.Message(1, "successfully list")
	ApiPath := os.Getenv("Api_path") + ":" + os.Getenv("PORT") + "/"
	whereCondition := ""
	if requestParams.Category != "" {
		whereCondition += " and lower(category) = '" + strings.ToLower(requestParams.Category) + "'"
	}
	if requestParams.Language != "" {
		whereCondition += " and language SIMILAR TO  '%(" + strings.Replace(requestParams.Language, ",", "|", 1) + ")%'"
	}
	objresponsetreding := make([]ResponseStruct, 0)
	if requestParams.Title == "" {
		querytrending := ""
		querytrending = "SELECT id, title, concat('" + ApiPath + "',video_url) as video_url,concat('" + ApiPath + "',thumb_url) as thumb_url, category, language, created_at, view_count FROM public.videostatus where 1=1 " + whereCondition + " order by view_count desc limit 20;"
		fmt.Println(querytrending)
		err := conn.Select(&objresponsetreding, querytrending)
		if err != nil {
			conn.Close()
			return u.Message(0, err.Error())
		}
	}
	whereConditionAll := ""
	if requestParams.CreatedOn != "" {
		whereConditionAll += " and created_at < " + requestParams.CreatedOn
	}
	objresponsedata := make([]ResponseStruct, 0)
	queryall := ""
	if requestParams.Title == "" {
		queryall = "SELECT id, title, concat('" + ApiPath + "',video_url) as video_url,concat('" + ApiPath + "',thumb_url) as thumb_url, category, language, created_at, view_count FROM public.videostatus where id not in (SELECT id FROM public.videostatus where 1=1 " + whereCondition + "  order by view_count desc limit 20) " + whereCondition + whereConditionAll + "  order by created_at desc limit 20;"
	} else {
		whereConditionAll += " and lower(title) like '%" + strings.ToLower(requestParams.Title) + "%'"
		queryall = "SELECT id, title, concat('" + ApiPath + "',video_url) as video_url,concat('" + ApiPath + "',thumb_url) as thumb_url, category, language, created_at, view_count FROM public.videostatus where 1=1 " + whereCondition + whereConditionAll + "  order by created_at desc limit 20;"
	}
	fmt.Println(queryall)
	err := conn.Select(&objresponsedata, queryall)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	queryalltogether := ""
	objresponsedataall := make([]ResponseStruct, 0)
	if requestParams.Title == "" {
		queryalltogether = "SELECT id, title, concat('" + ApiPath + "',video_url) as video_url,concat('" + ApiPath + "',thumb_url) as thumb_url, category, language, created_at, view_count FROM public.videostatus where 1=1 " + whereCondition + whereConditionAll + "  order by created_at desc;"
	} else {
		whereConditionAll += " and lower(title) like '%" + strings.ToLower(requestParams.Title) + "%'"
		queryalltogether = "SELECT id, title, concat('" + ApiPath + "',video_url) as video_url,concat('" + ApiPath + "',thumb_url) as thumb_url, category, language, created_at, view_count FROM public.videostatus where 1=1 " + whereCondition + whereConditionAll + "  order by created_at desc;"
	}
	fmt.Println(queryalltogether)
	err = conn.Select(&objresponsedataall, queryalltogether)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	data["trending"] = objresponsetreding
	data["all"] = objresponsedata
	data["alltogether"] = objresponsedataall
	resp["response_data"] = data
	conn.Close()
	return resp
}

func ViewCount(requestParams *RequestViewCountDelete) map[string]interface{} {
	conn := PgCon()
	resp := u.Message(1, "update successfully")
	_, err := conn.Exec("update videostatus set view_count = view_count +1 where id = '" + strconv.Itoa(requestParams.ID) + "'::int;")
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	conn.Close()
	return resp
}
func DeleteVideo(requestParams *RequestViewCountDelete) map[string]interface{} {
	conn := PgCon()
	video_url, thumb_url := "", ""
	resp := u.Message(1, "delete successfully")
	err := conn.QueryRow("SELECT video_url, thumb_url FROM public.videostatus where id = '"+strconv.Itoa(requestParams.ID)+"'::int;").Scan(&video_url, &thumb_url)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	_, err = conn.Exec("delete from videostatus where id = '" + strconv.Itoa(requestParams.ID) + "'::int;")
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	go DeleteFile(video_url)
	go DeleteFile(thumb_url)
	conn.Close()
	return resp
}
func DeleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if err != nil {
		fmt.Println("deelete file", err.Error())
	}

	fmt.Println("==> done deleting file")
}
