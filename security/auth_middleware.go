package security

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Ayikoandrew/server/database"
	"github.com/Ayikoandrew/server/types"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateAccessTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		access := os.Getenv("ACCESS_TOKEN")
		access_token, err := r.Cookie("access-token")

		if err != nil {
			slog.Error("Invalid access token")
		}

		token, err := jwt.ParseWithClaims(access_token.Value, &types.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %+v", t.Header["alg"])
			}
			return []byte(access), nil
		})

		if err != nil {
			slog.Error("Invalid access token")
		}

		claim, ok := token.Claims.(types.CustomClaims)

		if !ok {
			slog.Error("Invalid access token")
		}

		userID := claim.Subject
		Token := database.Get(userID, context.Background())
		redisToken, err := Token.Result()
		if err != nil {
			slog.Error("Invalid access token")
		}

		if access != redisToken {
			slog.Error("Invalid access token")
		}
		next(w, r)
	}
}
