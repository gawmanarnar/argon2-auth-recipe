package main

import (
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
}

// Login - http handler for user login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

// Register - http handler for user registration
func Register(w http.ResponseWriter, r *http.Request) {
	credentials := &LoginCredentials{}
	if err := json.NewDecoder(r.Body).Decode(credentials); err != nil {
		respondWithError(w, http.StatusBadRequest, "")
		return
	}

	hash, salt, err := crypto.Hash(credentials.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	_, err = s.DB.Query("INSERT INTO users (email, password, salt) VALUES ($1, $2, $3)", credentials.Username, hash, salt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	fmt.Println("User created")
}
