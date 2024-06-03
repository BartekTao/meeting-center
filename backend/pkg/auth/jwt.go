package auth

import (
	"log"
	"os"
	"time"

	"github.com/BartekTao/nycu-meeting-room-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
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

func (j *jwtHandler) GenerateJWT(userinfo *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   *userinfo.ID,
		"name":  userinfo.Name,
		"email": userinfo.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString(j.jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
