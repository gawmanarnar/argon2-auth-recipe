package crypto

import (
	"bytes"
	"encoding/base64"
	"errors"

	"github.com/gorilla/securecookie"
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

// doHash - actually does the hashing
func doHash(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, timeCost, memoryCost, threads, keyLength)
}

// GenerateRandomKey - wraps gorilla/securecookie's GenerateRandomKey
func GenerateRandomKey(length int) []byte {
	return securecookie.GenerateRandomKey(length)
}

// Hash - generates the Argon2i hash of a given password
// returns the hash and the salt that was used to create the hash
// these values are base64 encoded
func Hash(password string) (string, string, error) {
	unencodedSalt := GenerateRandomKey(32)
	if unencodedSalt == nil {
		return "", "", errors.New("Hash failed")
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
