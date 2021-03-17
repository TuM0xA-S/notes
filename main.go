package main

import (
	"log"
	"net/http"
	"notes/controllers"
	"notes/models"

	. "notes/config"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GetRouter returns prepared router
func GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.HandleFunc("/api/user/create", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Login).Methods("POST")
	router.HandleFunc("/api/me/notes", controllers.GetNotes).Methods("GET")
	router.HandleFunc("/api/me/notes/create", controllers.CreateNote).Methods("POST")
	router.HandleFunc("/api/me/notes/{note_id:[0-9]+}", controllers.NoteDetails).Methods("GET")
	router.HandleFunc("/api/me/notes/{note_id:[0-9]+}/remove", controllers.NoteRemove).Methods("POST")
	router.HandleFunc("/api/me", controllers.UserDetails).Methods("GET")
	return router
}

func jsonMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func main() {
	conn, err := gorm.Open(mysql.Open(Cfg.DBURI))
	if err != nil {
		log.Fatal("when connecting to db:", err)
	}
	models.Init(conn)
	models.Migrate()

	router := GetRouter()

	if err := http.ListenAndServe(Cfg.Host, router); err != nil {
		log.Printf("when serving: %v", err)
	}
}
