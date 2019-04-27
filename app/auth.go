package app

import (
	"fmt"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
)

var Authentication = func(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(0, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}
		conn := models.PgCon()
		id := ""
		err := conn.QueryRow("select id::character varying from user_master where auth_token=$1", tokenHeader).Scan(&id)
		if err != nil {
			fmt.Println("Login error:------")
			response = u.Message(0, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Sprintf("User %", id) //Useful for monitoring
		r.Header.Add("user_id", id)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
