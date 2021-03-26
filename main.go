package main

import (
	"fmt"
	"log"
	"net/http"
	"notes/controllers"
	"notes/models"
	"notes/util"
	"os"

	. "notes/config"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//FormatPanicError just responds with json
func (i InternalServerErrorResponder) FormatPanicError(rw http.ResponseWriter, _ *http.Request, _ *negroni.PanicInformation) {
	util.RespondWithError(rw, 500, "server internal error")
}

//InternalServerErrorResponder [lol i hate suppressing that warning messages]
type InternalServerErrorResponder struct{}

// GetRouter returns prepared router
func GetRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/notes", controllers.PublishedNotesList).Methods("GET")
	router.HandleFunc("/api/users", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/me", controllers.Login).Methods("POST")
	router.HandleFunc("/api/me/notes", controllers.NotesList).Methods("GET")
	router.HandleFunc("/api/me/notes", controllers.CreateNote).Methods("POST")
	router.HandleFunc("/api/me/notes/{note_id:[0-9]+}", controllers.NoteDetails).Methods("GET")
	router.HandleFunc("/api/me/notes/{note_id:[0-9]+}", controllers.NoteRemove).Methods("DELETE")
	router.HandleFunc("/api/me/notes/{note_id:[0-9]+}", controllers.NoteUpdate).Methods("PUT")
	router.HandleFunc("/api/me", controllers.UserDetails).Methods("GET")
	router.HandleFunc("/api/notes/{note_id:[0-9]+}", controllers.PublishedNoteDetail).Methods("GET")
	router.HandleFunc("/api/users/{user_id:[0-9]+}", controllers.AnotherUserDetail).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(controllers.NotFound)

	n := negroni.New()

	logger := log.New(os.Stdout, "[notes]", 0)

	loggerMid := negroni.NewLogger()
	loggerMid.ALogger = logger
	n.Use(loggerMid)

	recoverMid := negroni.NewRecovery()
	recoverMid.Logger = logger
	recoverMid.Formatter = InternalServerErrorResponder{}
	n.Use(recoverMid)

	n.UseHandler(router)

	return n
}

func main() {
	DBURI := fmt.Sprintf("root:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", Cfg.DBPassword, Cfg.DBHost, Cfg.DBName)
	log.Println("db_uri:", DBURI)
	conn, err := gorm.Open(mysql.Open(DBURI), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
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
