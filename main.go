package main

import (
	"fmt"
	"go-contacts/app"
	"go-contacts/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	STATIC_DIR := "/" + os.Getenv("video_path")
	router := mux.NewRouter()

	router.HandleFunc("/api/upload", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/list", controllers.VideoList).Methods("POST")
	router.HandleFunc("/api/viewcount", controllers.ViewCount).Methods("POST")
	router.HandleFunc("/api/deletevideo", controllers.DeleteVideo).Methods("POST")

	// -------------------------- Instagram clone -------------------------------------
	router.HandleFunc("/api/user_signup", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user_signin", controllers.CheckMobileOTP).Methods("POST")
	router.HandleFunc("/api/profile_update", app.Authentication(controllers.UpdateProfile)).Methods("POST")
	router.HandleFunc("/api/check_user_name_exists_or_not", app.Authentication(controllers.CheckUserName)).Methods("POST")
	router.HandleFunc("/api/upload_post", app.Authentication(controllers.UploadPost)).Methods("POST")
	router.HandleFunc("/api/following_request", app.Authentication(controllers.FollowingRequest)).Methods("POST")
	router.HandleFunc("/api/follow_get", app.Authentication(controllers.FollowGet)).Methods("POST")
	router.HandleFunc("/api/followers_action", app.Authentication(controllers.FollowersAction)).Methods("POST")
	router.HandleFunc("/api/get_user_list", app.Authentication(controllers.GetUserList)).Methods("POST")
	router.HandleFunc("/api/get_post_list", app.Authentication(controllers.GetListOfPost)).Methods("POST")

	router.
		PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
