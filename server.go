package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	"github.com/gimmeasandwich/argon2-auth-recipe/crypto"
	"github.com/gimmeasandwich/argon2-auth-recipe/middleware"
	"github.com/gimmeasandwich/argon2-auth-recipe/views"
	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	_ "github.com/lib/pq"
)

// WebServer - The web server
type WebServer struct {
	Router   *chi.Mux
	DB       *sql.DB
	Sessions *scs.Manager
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
	sessionManager := scs.NewCookieManager(string(crypto.GenerateRandomKey(32)))
	sessionManager.Lifetime(time.Hour * 24 * 30) // One month
	sessionManager.Persist(true)

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

	log.Fatal(http.ListenAndServe(":3000", middleware.SecureHeaders(middleware.Logger(csrfMiddleware(s.Router)))))
}

// LoggedIn - helper function to determine if a user is logged in
func (s *WebServer) LoggedIn(r *http.Request) (bool, error) {
	session := s.Sessions.Load(r)
	loggedIn, err := session.Exists("UserId")

	if err != nil {
		return false, err
	}

	return loggedIn, nil
}

// RequireLogin - wraps a route to require the user to login
func (s *WebServer) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggedIn, err := s.LoggedIn(r)

		if err != nil {
			// TODO: render error page
			return
		}

		if !loggedIn {
			http.Redirect(w, r, "/login", 302)
			return
		}

		next.ServeHTTP(w, r)
	})
}
