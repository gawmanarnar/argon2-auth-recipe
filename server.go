package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

// WebServer - The web server
type WebServer struct {
	Router *chi.Mux
	DB     *sql.DB
}

// Init - Initalizes our router and database connection
func (s *WebServer) Init(user, password, dbname string) {
	s.Router = chi.NewRouter()

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	var err error
	s.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

// Run - starts the server
func (s *WebServer) Run() {
	log.Fatal(http.ListenAndServe(":3000", s.Router))
}
