package config

import (
	"github.com/gimmeasandwich/argon2-auth-recipe/crypto"
	"github.com/spf13/viper"
)

// DB - contains information to connect to the database
var DB PostgresConfig

// Secrets - contains secret configuration
var Secrets SecretConfig

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	postgres := viper.Sub("postgres")
	err = postgres.Unmarshal(&DB)
	if err != nil {
		panic(err)
	}

	var forceWrite = false

	// Get or generate csrf key
	csrf := viper.GetString("secrets.csrf")
	if csrf == "" {
		csrf = crypto.ToBase64(crypto.GenerateRandomKey(32))
		viper.Set("secrets.csrf", csrf)
		forceWrite = true
	}
	Secrets.Csrf = csrf

	// Get or generate key to sign cookies
	cookie := viper.GetString("secrets.cookie")
	if cookie == "" {
		cookie = crypto.ToBase64(crypto.GenerateRandomKey(32))
		viper.Set("secrets.cookie", cookie)
		forceWrite = true
	}
	Secrets.Cookie = cookie

	if forceWrite {
		viper.WriteConfig()
	}
}

// PostgresConfig - contains information necessary to create a dbconnection
type PostgresConfig struct {
	Username string
	Password string
	Name     string
}

// SecretConfig - contains secret configuration
type SecretConfig struct {
	Csrf   string
	Cookie string
}

// GetDecodedCsrf - Gets the decoded csrf token
func (s *SecretConfig) GetDecodedCsrf() []byte {
	return GetDecoded(s.Csrf)
}

// GetDecodedCookie - Gets the decoded cookie token
func (s *SecretConfig) GetDecodedCookie() []byte {
	return GetDecoded(s.Cookie)
}

// GetDecoded - take a base64 encoded string and returns the decoded byte array
func GetDecoded(encoded string) []byte {
	decoded, err := crypto.FromBase64(encoded)
	if err != nil {
		panic(err)
	}
	return decoded
}
