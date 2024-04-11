package auth

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtHandler struct {
	jwtKey []byte
}

func NewJWTHandler() *jwtHandler {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Fatal("You must set the JWT_KEY environment variable")
	}
	return &jwtHandler{jwtKey: []byte(jwtKey)}
}

func (j *jwtHandler) GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(j.jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
