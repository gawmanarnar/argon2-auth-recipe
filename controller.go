package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gimmeasandwich/argon2-auth-recipe/crypto"
	"github.com/gorilla/schema"
	_ "github.com/lib/pq"
)

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
}

// LoginCredentials - struct containing the users login credentials
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Salt     string `json:"salt,omitempty"`
}

func (c *LoginCredentials) getCredentials(db *sql.DB) error {
	row := db.QueryRow("SELECT email, password, salt FROM users WHERE email=$1", c.Email)
	return row.Scan(&c.Email, &c.Password, &c.Salt)
}

// Login - http handler for user login
func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	credentials := &LoginCredentials{}
	if err := decoder.Decode(credentials, r.PostForm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storedCredentials := LoginCredentials{Email: credentials.Email}
	if err := storedCredentials.getCredentials(s.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if crypto.VerifyHash(credentials.Password, storedCredentials.Password, storedCredentials.Salt) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// Register - http handler for user registration
func Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	credentials := &LoginCredentials{}
	if err := decoder.Decode(credentials, r.PostForm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, salt, err := crypto.Hash(credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = s.DB.Query("INSERT INTO users (email, password, salt) VALUES ($1, $2, $3)", credentials.Email, hash, salt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("User created")
}
