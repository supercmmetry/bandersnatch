package middleware

import (
	"bandersnatch/api"
	"context"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
)

type Token struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

type JwtContextKey string

func JwtAuth(next http.Handler) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			api.RespWrap(w, http.StatusForbidden, "auth token missing")
			return
		}

		tk := &Token{}
		token, err := jwt.ParseWithClaims(tokenHeader, tk, func (t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_PASSWORD")), nil
		})

		if err != nil {
			api.RespWrap(w, http.StatusForbidden, "malformed auth token")
			return
		}

		if !token.Valid {
			api.RespWrap(w, http.StatusForbidden, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), JwtContextKey("token"), tk)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
