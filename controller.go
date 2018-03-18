package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gimmeasandwich/argon2-auth-recipe/crypto"
	_ "github.com/lib/pq"
)

// LoginCredentials - struct containing the users login credentials
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt,omitempty"`
}

func (c *LoginCredentials) getCredentials(db *sql.DB) error {
	row := db.QueryRow("SELECT email, password, salt FROM users WHERE email=$1", c.Username)
	return row.Scan(&c.Username, &c.Password, &c.Salt)
}

// Login - http handler for user login
func Login(w http.ResponseWriter, r *http.Request) {
	credentials := &LoginCredentials{}
	err := json.NewDecoder(r.Body).Decode(credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storedCredentials := LoginCredentials{Username: credentials.Username}
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
	credentials := &LoginCredentials{}
	if err := json.NewDecoder(r.Body).Decode(credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, salt, err := crypto.Hash(credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = s.DB.Query("INSERT INTO users (email, password, salt) VALUES ($1, $2, $3)", credentials.Username, hash, salt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("User created")
}
