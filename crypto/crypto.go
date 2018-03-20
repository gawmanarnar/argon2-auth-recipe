package crypto

import (
	"bytes"
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

// doHash - actually does the hashing
func doHash(password, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLength)
}

// Hash - generates the Argon2i hash of a given password
// returns the hash and the salt that was used to create the hash
// these values are base64 encoded
func Hash(password string) (string, string, error) {
	unencodedSalt, err := generateSalt()
	if err != nil {
		return "", "", err
	}

	unencodedHash := doHash([]byte(password), unencodedSalt)
	if didHashFail(unencodedHash) {
		return "", "", errors.New("Hash failed")
	}

	hash := base64.StdEncoding.EncodeToString(unencodedHash)
	salt := base64.StdEncoding.EncodeToString(unencodedSalt)

	return hash, salt, nil
}

// VerifyHash - takes a password, a base64 encoded hash, and a base64 encoded salt
// returns true if the password matches the hash
func VerifyHash(password, hash, salt string) bool {
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}

	decodedSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false
	}

	testHash := doHash([]byte(password), decodedSalt)
	return bytes.Equal(decodedHash, testHash)
}
