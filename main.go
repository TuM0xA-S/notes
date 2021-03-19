package main

import (
	"fmt"
	"log"
	"net/http"
	"notes/controllers"
	"notes/models"

	. "notes/config"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GetRouter returns prepared router
func GetRouter() http.Handler {
	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.HandleFunc("/api/notes", controllers.PublishedNotesList).Methods("GET")
	router.HandleFunc("/api/user/create", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Login).Methods("POST")
	router.HandleFunc("/api/me/notes", controllers.NotesList).Methods("GET")
	router.HandleFunc("/api/me/notes/create", controllers.CreateNote).Methods("POST")
	router.HandleFunc("/api/me/notes/{note_id:[0-9]+}", controllers.NoteDetails).Methods("GET")
	router.HandleFunc("/api/me/notes/{note_id:[0-9]+}", controllers.NoteRemove).Methods("DELETE")
	router.HandleFunc("/api/me/notes/{note_id:[0-9]+}", controllers.NoteUpdate).Methods("PUT")
	router.HandleFunc("/api/me", controllers.UserDetails).Methods("GET")

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(router)
	return n
}

func jsonMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func main() {
	DBURI := fmt.Sprintf("root:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", Cfg.DBPassword, Cfg.DBHost, Cfg.DBName)
	log.Println("db_uri:", DBURI)
	conn, err := gorm.Open(mysql.Open(DBURI))
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
