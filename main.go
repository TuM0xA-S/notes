package main

import (
	"contacts/controllers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("when loading env: %v", err)
	}
}

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Login).Methods("POST")
	router.HandleFunc("/api/me/contacts", controllers.GetContacts).Methods("GET")
	router.HandleFunc("/api/me/contacts/create", controllers.CreateContact).Methods("POST")
	return router
}

func main() {

	router := getRouter()
	host := os.Getenv("port")
	if host == "" {
		host = ":8000"
	}

	if err := http.ListenAndServe(host, router); err != nil {
		log.Printf("when serving: %v", err)
	}
}
