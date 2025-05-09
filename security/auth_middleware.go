package security

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddlewareJwt(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("auth-x-token")
		token, err := ValidateAccessToken(tokenStr)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, token)
			return
		}
		next(w, r)
	}
}

func ValidateAccessToken(tokenStr string) (*jwt.Token, error) {
	access := os.Getenv("ACCESS_TOKEN")
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %+v", t.Header["alg"])
		}
		return []byte(access), nil
	})
}

func ValidateRefreshToken(tokenStr string) (*jwt.Token, error) {
	refresh := os.Getenv("REFRESH_TOKEN")
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %+v", t.Header["alg"])
		}
		return []byte(refresh), nil
	})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
