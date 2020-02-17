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
	token []byte
}

func Hostname() string {
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
	return time.Now().Add(2 * time.Hour).Unix()
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
		token: []byte(tokenString),
	}
}

func Validate() {
}
