package models

import (
	"database/sql"
	"fmt"
	u "go-contacts/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/chr4/pwgen"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
	"gopkg.in/guregu/null.v3/zero"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

// //a struct to rep user account
type Account struct {
	Mobile            string `json:"mobile_no" validate:"required"`
	CountryCode       string `json:"country_code" validate:"required"`
	VerificationToken string `json:"verification_token"`
	AuthToken         string `json:"auth_token"`
}

type UserData struct {
	ID             int         `json:"id" db:"id"`
	FirstName      zero.String `json:"first_name" db:"first_name"`
	LastName       zero.String `json:"last_name" db:"last_name"`
	AuthToken      string      `json:"auth_token" db:"auth_token"`
	ProfilePicture zero.String `json:"profile_picture" db:"profile_picture"`
	UserName       zero.String `json:"username" db:"username"`
	ProfileStatus  zero.String `json:"profile_status" db:"profile_status"`
	LocalID        string      `json:"local_id" db:"local_id"`
}

type RequestUserUpdate struct {
	ID            int                   `json:"id" db:"id"`
	FirstName     string                `json:"first_name" validate:"required"`
	LastName      string                `json:"last_name" validate:"required"`
	File          multipart.File        `json:"file"`
	FileHeader    *multipart.FileHeader `json:"fileheader"`
	UserName      string                `json:"username" validate:"required"`
	ProfileStatus string                `json:"profile_status"`
}

type UserName struct {
	Username string `json:"username" validate:"required"`
}

func CreateAccount(requestParams *Account) map[string]interface{} {
	resp := make(map[string]interface{})
	// data := make(map[string]interface{})
	conn := PgCon()
	guid := xid.New()
	requestParams.AuthToken = pwgen.AlphaNum(32)
	if !u.RowExists(conn, "select id from user_master where country_code=$1 and mobile_no=$2", requestParams.CountryCode, requestParams.Mobile) {
		_, err := conn.Exec("insert into user_master(country_code, mobile_no, auth_token, local_id) values($1,$2,$3,$4)", requestParams.CountryCode, requestParams.Mobile, requestParams.AuthToken, guid.String())
		if err != nil {
			conn.Close()
			return u.Message(0, err.Error())
		}
		resp = u.Message(1, "User added successfully. Please check your mobile for OTP.")
		resp["auth_token"] = requestParams.AuthToken
	} else {
		resp = u.Message(0, "User already exists please login in.")
	}
	conn.Close()
	return resp
}

func CheckMobileOTP(requestParams *Account) map[string]interface{} {
	resp := make(map[string]interface{})
	// data := make(map[string]interface{})
	conn := PgCon()
	userData := UserData{}
	err := conn.QueryRowx("select id,first_name,last_name,auth_token,profile_picture,username,profile_status,local_id from user_master where country_code=$1 and mobile_no=$2", requestParams.CountryCode, requestParams.Mobile).StructScan(&userData)
	if err != nil && err != sql.ErrNoRows {
		conn.Close()
		return u.Message(0, err.Error())
	} else if err == sql.ErrNoRows {
		resp = u.Message(0, "Invalid credentials. Please check your mobile number or sign up.")
	} else {
		requestParams.AuthToken = pwgen.AlphaNum(32)
		_, err = conn.Exec("update user_master set auth_token = $1,updated_at = (date_part('epoch'::text, now()) * (1000)::double precision) where id=$2", requestParams.AuthToken, userData.ID)
		if err != nil {
			conn.Close()
			return u.Message(0, err.Error())
		}
		response := make(map[string]interface{})
		response["auth_token"] = requestParams.AuthToken
		response["first_name"] = userData.FirstName.String
		response["last_name"] = userData.LastName.String
		response["username"] = userData.UserName.String
		response["profile_status"] = userData.ProfileStatus.String
		if userData.ProfilePicture.String != "" {
			response["profile_picture"] = os.Getenv("Api_path") + ":" + os.Getenv("PORT") + "/" + os.Getenv("user_profile") + userData.LocalID + "/" + userData.ProfilePicture.String
		} else {
			response["profile_picture"] = userData.ProfilePicture.String
		}
		resp = u.Message(1, "login successfully.")
		resp["data"] = response
	}
	conn.Close()
	return resp
}

func UpdateProfile(requestParams *RequestUserUpdate) map[string]interface{} {
	resp := make(map[string]interface{})
	conn := PgCon()
	fileName := ""
	img, local_id := "", ""
	if u.RowExists(conn, "select id from user_master where username=$2 and id != $1", requestParams.ID, requestParams.UserName) {
		conn.Close()
		return u.Message(0, "Please use different unique name.")
	}
	err := conn.QueryRow("select profile_picture,local_id from user_master where id = $1", requestParams.ID).Scan(&img, &local_id)
	if err != nil && err != sql.ErrNoRows {
		conn.Close()
		return u.Message(0, err.Error())
	}
	if requestParams.FileHeader != nil {
		fileName = "user_profile" + filepath.Ext(requestParams.FileHeader.Filename)
		if img != "" {
			go DeleteFile(img)
		}
		go SaveFileUser(requestParams.File, fileName, "./"+os.Getenv("user_profile")+local_id+"/")
	} else if requestParams.FileHeader == nil {
		fileName = img
	}
	_, err = conn.Exec("update user_master set first_name=$1,last_name=$2,profile_picture=$3,username=$5,profile_status=$6,updated_at = (date_part('epoch'::text, now()) * (1000)::double precision) where id=$4", requestParams.FirstName, requestParams.LastName, fileName, requestParams.ID, requestParams.UserName, requestParams.ProfileStatus)
	if err != nil {
		conn.Close()
		return u.Message(0, err.Error())
	}
	resp = u.Message(1, "Profile Updated successfully")
	conn.Close()
	return resp
}

func SaveFileUser(File multipart.File, FileName string, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	f, err := os.OpenFile(path+FileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("error :- ", err)
	}
	defer f.Close()
	io.Copy(f, File)
	fmt.Println("successfully add")
}

func CheckUserName(requestParams *UserName) map[string]interface{} {
	resp := make(map[string]interface{})
	// data := make(map[string]interface{})
	conn := PgCon()
	if u.RowExists(conn, "select id from user_master where username=$1", requestParams.Username) {
		resp = u.Message(0, "Please use different unique name.")
	} else {
		resp = u.Message(1, "you can use this name.")
	}
	conn.Close()
	return resp
}

// func Login(email, password string) map[string]interface{} {

// 	account := &Account{}
// 	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return u.Message(false, "Email address not found")
// 		}
// 		return u.Message(false, "Connection error. Please retry")
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
// 	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
// 		return u.Message(false, "Invalid login credentials. Please try again")
// 	}
// 	//Worked! Logged In
// 	account.Password = ""

// 	//Create JWT token
// 	tk := &Token{UserId: account.ID}
// 	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
// 	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
// 	account.Token = tokenString //Store the token in the response

// 	resp := u.Message(true, "Logged In")
// 	resp["account"] = account
// 	return resp
// }

// func GetUser(u uint) *Account {

// 	acc := &Account{}
// 	GetDB().Table("accounts").Where("id = ?", u).First(acc)
// 	if acc.Email == "" { //User not found!
// 		return nil
// 	}

// 	acc.Password = ""
// 	return acc
// }
