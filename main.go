package main

import (
	"fmt"
	"net/http"
	"os"

	"./app"
	"./controllers"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/user/verify/{uid}/{token}", controllers.VerifyEMail).Methods("GET")
	router.HandleFunc("/api/sections", controllers.GetSections).Methods("GET")
	router.HandleFunc("/api/posts", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/api/posts", controllers.GetPosts).Methods("GET")
	router.HandleFunc("/api/posts/{id}", controllers.GetPost).Methods("GET")
	router.HandleFunc("/api/posts/{id}", controllers.DeletePost).Methods("DELETE")
	router.HandleFunc("/api/comments", controllers.CreateComment).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	fmt.Println("Run in " + port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
