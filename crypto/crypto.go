package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	timeCost   = 3
	memoryCost = 32 * 1024
	threads    = 4
	keyLength  = 64
)

// didHashFail - checks to see if the hash is empty (all zeros)
func didHashFail(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

// generateSalt - generates a 32 byte salt
func generateSalt() ([]byte, error) {
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

// Hash - generates the Argon2i hash of a given password
// returns the hash and the salt that was used to create the hash
// these values are base64 encoded
func Hash(password string) (string, string, error) {
	unencodedSalt, err := generateSalt()
	if err != nil {
		return "", "", err
	}

	unencodedHash := argon2.Key([]byte(password), unencodedSalt, timeCost, memoryCost, threads, keyLength)
	if didHashFail(unencodedHash) {
		return "", "", errors.New("Hash failed")
	}

	hash := base64.StdEncoding.EncodeToString(unencodedHash)
	salt := base64.StdEncoding.EncodeToString(unencodedSalt)

	return hash, salt, nil
}
