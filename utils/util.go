package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	validator "gopkg.in/go-playground/validator.v9"
)

func Message(status int, message string) map[string]interface{} {
	return map[string]interface{}{"response_code": status, "response_message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ValidateObject(requestBody interface{}) (map[string]interface{}, bool) {
	validate := validator.New()
	fmt.Println(requestBody)
	err := validate.Struct(requestBody)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return map[string]interface{}{"response_code": 0, "response_message": err.Error()}, false
		}
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
			return map[string]interface{}{"response_code": 0, "response_message": err.Field() + " field is required"}, false
		}
	}
	return nil, true
}

func RowExists(db *sqlx.DB, query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("error checking if row exists '%s' %v", args, err)
	}
	return exists
}
