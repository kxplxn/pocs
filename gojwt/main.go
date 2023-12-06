package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	key := []byte("asdfjklhaskdfjahsdkl")

	tkRaw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     "bob123",
		"isAdmin": true,
		"teamID":  "asdfasdfads",
		"exp":     time.Now().Add(1 * time.Hour).UTC().Unix(),
	}).SignedString(key)
	if err != nil {
		log.Fatalln("failed to generate token:", err)
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(
		tkRaw, &claims, func(token *jwt.Token) (any, error) {
			return key, nil
		},
	)
	if err != nil {
		log.Fatalln("failed to parse token:", err)
	}

	fmt.Printf("claims: %+v\n", claims)
}
