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
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		conn := models.PgCon()
		id, local_id := "", ""
		err := conn.QueryRow("select id::character varying,local_id from user_master where auth_token=$1", tokenHeader).Scan(&id, &local_id)
		if err != nil {
			conn.Close()
			fmt.Println("Login error:------")
			response = u.Message(0, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		conn.Close()
		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Printf("User %", id) //Useful for monitoring
		r.Header.Add("user_id", id)
		r.Header.Add("user_local_id", local_id)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
