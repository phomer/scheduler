package accounts

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/phomer/scheduler/log"

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

func CreateToken() *Token {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 15000,
		Issuer:    "scheduler",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SecretKey())
	if err != nil {
		fmt.Println("Failed to Sign Token ", err, token)
		panic("Goodbye")
	}

	// Test it
	_, err = jwt.ParseWithClaims(tokenString, claims, GetKey)
	if err != nil {
		log.Fatal("Tokens are failing", err)
	}

	return &Token{
		Signed: tokenString,
	}
}

func NewToken(tokenString string) *Token {
	return &Token{
		Signed: tokenString,
	}
}

func Validate(token *Token) bool {

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 15000,
		Issuer:    "scheduler",
	}

	tokenStatus, err := jwt.ParseWithClaims(token.Signed, claims, GetKey)
	if err != nil {
		log.Fatal("ParseWithClaims", err, tokenStatus, token.Signed)
	}

	return tokenStatus.Valid
}

func GetKey(token *jwt.Token) (interface{}, error) {
	// TODO: Interrupted by cat needing food :-(, will fix later.
	return SecretKey(), nil
}
