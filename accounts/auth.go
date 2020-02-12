package accounts

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	Signed string
}

func Hostname() string {
	// TODO: Replace this, if it isn't commonly set on servers?
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Finding Hostname", err)
	}
	return hostname
}

func Username() string {
	entry, err := user.Current()
	if err != nil {
		log.Fatal("Finding Username", err)
	}
	return entry.Username
}

// # of Seconds since ...
func ExpiryDate() int64 {

	// TODO: Implement some form of reissing the token, besides reregistering
	return time.Now().Add(24 * 60 * time.Hour).Unix()
}

func SecretKey() []byte {
	return []byte("RandomizeThisPlease")
}

func NewToken() *Token {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["admin"] = false      // TODO: Check this?
	claims["name"] = Username()  // TODO: Check this?
	claims["exp"] = ExpiryDate() // TODO: Check this?

	tokenString, err := token.SignedString(SecretKey())
	if err != nil {
		fmt.Println("Failed to Sign Token", err, token)
		panic("Goodbye")
	}

	return &Token{
		Signed: tokenString,
	}
}

func Validate(token *Token) bool {

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * 60 * time.Hour).Unix(),
	}

	tokenString, err := jwt.ParseWithClaims(token.Signed, claims, GetKey)
	if err != nil {
		log.Fatal("ParseWithClaims", err)
	}

	return tokenString.Valid
}

func GetKey(token *jwt.Token) (interface{}, error) {
	// TODO: Interrupted by cat needing food :-(, will fix later.
	return []byte("Secret Key"), nil
}
