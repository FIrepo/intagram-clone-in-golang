package models

import (
	"database/sql"
	u "go-contacts/utils"
	"mime/multipart"

	"github.com/chr4/pwgen"
	"gopkg.in/guregu/null.v3/zero"

	"github.com/dgrijalva/jwt-go"
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
	VerificationToken string `json:"verification_token" validate:"required"`
	AuthToken         string `json:"auth_token"`
}

type UserData struct {
	ID        int         `json:"id" db:"id"`
	FirstName zero.String `json:"first_name" db:"first_name"`
	LastName  zero.String `json:"last_name" db:"last_name"`
	AuthToken string      `json:"auth_token" db:"auth_token"`
}

type RequestUserUpdate struct {
	ID        int                   `json:"id" db:"id"`
	FirstName string                `json:"first_name"`
	LastName  string                `json:"last_name"`
	File      *multipart.FileHeader `json:"file" validate:"required"`
	UserName  string                `json:"username"`
}

func CreateAccount(requestParams *Account) map[string]interface{} {
	resp := make(map[string]interface{})
	// data := make(map[string]interface{})
	conn := PgCon()
	requestParams.AuthToken = pwgen.AlphaNum(32)
	if !u.RowExists(conn, "select id from user_master where country_code=$1 and mobile_no=$2", requestParams.CountryCode, requestParams.Mobile) {
		_, err := conn.Exec("insert into user_master(country_code, mobile_no, auth_token,verification_token) values($1,$2,$3,$4)", requestParams.CountryCode, requestParams.Mobile, requestParams.AuthToken, requestParams.VerificationToken)
		if err != nil {
			conn.Close()
			return u.Message(0, err.Error())
		}
		resp = u.Message(1, "User added successfully. Please check your mobile for OTP.")
	} else {
		_, err := conn.Exec("update user_master set  auth_token=$3,verification_token=$4,updated_at = (date_part('epoch'::text, now()) * (1000)::double precision) where country_code=$1 and mobile_no=$2", requestParams.CountryCode, requestParams.Mobile, requestParams.AuthToken, requestParams.VerificationToken)
		if err != nil {
			conn.Close()
			return u.Message(0, err.Error())
		}
		resp = u.Message(1, "Please check your mobile for OTP.")
	}
	conn.Close()
	return resp
}

func CheckMobileOTP(requestParams *Account) map[string]interface{} {
	resp := make(map[string]interface{})
	// data := make(map[string]interface{})
	conn := PgCon()
	userData := UserData{}
	err := conn.QueryRowx("select id,first_name,last_name,auth_token from user_master where country_code=$1 and mobile_no=$2 and verification_token=$3", requestParams.CountryCode, requestParams.Mobile, requestParams.VerificationToken).StructScan(&userData)
	if err != nil && err != sql.ErrNoRows {
		conn.Close()
		return u.Message(0, err.Error())
	} else if err == sql.ErrNoRows {
		resp = u.Message(0, "Mobile OTP is wrong.")
	} else {
		response := make(map[string]interface{})
		response["auth_token"] = userData.AuthToken
		response["first_name"] = userData.FirstName
		response["last_name"] = userData.LastName
		resp = u.Message(1, "OTP verification successfully.")
		resp["data"] = response
	}
	conn.Close()
	return resp
}

func UpdateProfile(requestParams *RequestUserUpdate) map[string]interface{} {
	resp := make(map[string]interface{})
	// data := make(map[string]interface{})
	// conn := PgCon()
	// userData := UserData{}
	// err := conn.QueryRowx("select id,first_name,last_name,auth_token from user_master where country_code=$1 and mobile_no=$2 and verification_token=$3", requestParams.CountryCode, requestParams.Mobile, requestParams.VerificationToken).StructScan(&userData)
	// if err != nil && err != sql.ErrNoRows {
	// 	conn.Close()
	// 	return u.Message(0, err.Error())
	// } else if err == sql.ErrNoRows {
	// 	resp = u.Message(0, "Mobile OTP is wrong.")
	// } else {
	// 	response := make(map[string]interface{})
	// 	response["auth_token"] = userData.AuthToken
	// 	response["first_name"] = userData.FirstName
	// 	response["last_name"] = userData.LastName
	// 	resp = u.Message(1, "OTP verification successfully.")
	// 	resp["data"] = response
	// }
	resp["data"] = requestParams
	// conn.Close()
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
