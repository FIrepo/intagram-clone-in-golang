package main

import (
	"fmt"
	"go-contacts/app"
	"go-contacts/controllers"
	u "go-contacts/utils"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	u.CreateFb()
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
	router.HandleFunc("/api/user_signup", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user_signin", controllers.CheckMobileOTP).Methods("POST")
	router.HandleFunc("/api/profile_update", app.Authentication(controllers.UpdateProfile)).Methods("POST")
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
