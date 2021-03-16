package main

import (
	"log"
	"net/http"
	"notes/controllers"
	"notes/models"

	"notes/config"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Cfg = config.Cfg

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Login).Methods("POST")
	router.HandleFunc("/api/me/notes", controllers.GetNotes).Methods("GET")
	router.HandleFunc("/api/me/notes/create", controllers.CreateNote).Methods("POST")
	return router
}

func main() {
	conn, err := gorm.Open(mysql.Open(Cfg.DBURI))
	if err != nil {
		log.Fatal("when connecting to db:", err)
	}
	models.Init(conn)
	models.Migrate()

	router := getRouter()

	if err := http.ListenAndServe(Cfg.Host, router); err != nil {
		log.Printf("when serving: %v", err)
	}
}
