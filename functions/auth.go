package api

import (
	"os"
	"time"

	"github.com/Ayikoandrew/server/types"
	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(accountId string) (string, error) {
	access := os.Getenv("ACCESS_TOKEN")

	claim := &types.CustomClaims{
		UserID: accountId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "liora-access",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   accountId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(access))
}

func CreateRefreshToken(accountId string) (string, error) {
	refresh := os.Getenv("REFRESH_TOKEN")

	claim := &types.CustomClaims{
		UserID: accountId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "liora-refresh",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   accountId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(refresh))
}
