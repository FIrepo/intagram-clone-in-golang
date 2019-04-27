package app

import (
	u "go-contacts/utils"
	"net/http"
)

var NotFoundHandler = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(0, "This resources was not found on our server"))
		next.ServeHTTP(w, r)
	})
}
