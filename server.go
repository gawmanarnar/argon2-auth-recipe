package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gimmeasandwich/argon2-auth-recipe/crypto"
	"github.com/gimmeasandwich/argon2-auth-recipe/views"
	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	_ "github.com/lib/pq"
)

// WebServer - The web server
type WebServer struct {
	Router *chi.Mux
	DB     *sql.DB
}

// SetupDB - Opens postgres database connection
func (s *WebServer) SetupDB(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	var err error
	s.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

// SetupRoutes - Initializes the routes for the application
func (s *WebServer) SetupRoutes() {
	s.Router = chi.NewRouter()
	s.Router.Post("/login", Login)
	s.Router.Post("/signup", Register)
	s.Router.Get("/login", views.RenderLogin)
	s.Router.Get("/signup", views.RenderSignup)
}

// Run - starts the server
func (s *WebServer) Run() {

	// Setup csrf protection
	csrfMiddleware := csrf.Protect(crypto.GenerateRandomKey(32), csrf.Secure(false))

	log.Fatal(http.ListenAndServe(":3000", csrfMiddleware(s.Router)))
}
