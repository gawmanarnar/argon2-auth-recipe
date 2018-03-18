package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/argon2"
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
	err := json.NewDecoder(r.Body).Decode(credentials)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "")
		return
	}

	unencodedSalt, err := generateSalt()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	unencodedHash := argon2.Key([]byte(credentials.Password), unencodedSalt, 3, 32*1024, 4, 64)
	if isByteArrayEmpty(unencodedHash) {
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	hash := base64.StdEncoding.EncodeToString(unencodedHash)
	salt := base64.StdEncoding.EncodeToString(unencodedSalt)

	fmt.Println(credentials.Username)
	fmt.Println(hash)
	fmt.Println(salt)
}
