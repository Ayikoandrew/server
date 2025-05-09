package security

import (
	"os"
	"time"

	"github.com/Ayikoandrew/server/types"
	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(account *types.User) (string, error) {
	access := os.Getenv("ACCESS_TOKEN")

	claim := &jwt.RegisteredClaims{
		Issuer:    "liora-access",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		Subject:   account.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(access))
}

func CreateRefreshToken(account *types.User) (string, error) {
	refresh := os.Getenv("REFRESH_TOKEN")

	claim := &jwt.RegisteredClaims{
		Issuer:    "liora-refresh",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		Subject:   account.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(refresh))
}
