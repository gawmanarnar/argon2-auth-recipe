package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/argon2"
)

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Generates a 32 byte salt
func generateSalt() ([]byte, error) {
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
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
		if allZero(unencodedHash) {
			respondWithError(w, http.StatusInternalServerError, "")
			return
		}

		hash := base64.StdEncoding.EncodeToString(unencodedHash)
		salt := base64.StdEncoding.EncodeToString(unencodedSalt)

		fmt.Println(credentials.Username)
		fmt.Println(hash)
		fmt.Println(salt)
	})

	log.Fatal(http.ListenAndServe(":3000", r))
}
