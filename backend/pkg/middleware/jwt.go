package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserCtxKey = contextKey("user")

type MeetingCenterClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")

			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := header[7:]

			token, err := jwt.ParseWithClaims(tokenStr, &MeetingCenterClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				log.Panicln(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else if claims, ok := token.Claims.(*MeetingCenterClaims); ok {
				ctx := context.WithValue(r.Context(), UserCtxKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			} else {
				log.Panicln("unknown claims type, cannot proceed")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}
