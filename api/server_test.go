package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/Ayikoandrew/server/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRefreshToken(t *testing.T) {
	testSecret := "test-secret-key-for-jwt-signing"

	server := &Server{}

	t.Run("Valid Token", func(t *testing.T) {
		claims := &types.CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "test-user-id",
				Issuer:    "liora-refresh",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testSecret))
		require.NoError(t, err)

		validatedClaims, err := server.validateRefreshToken(tokenString, testSecret)

		assert.NoError(t, err)
		assert.NotNil(t, validatedClaims)
		assert.Equal(t, "test-user-id", validatedClaims.Subject)
		assert.Equal(t, "liora-refresh", validatedClaims.Issuer)
	})

	t.Run("Invalid Secret", func(t *testing.T) {
		claims := &types.CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "test-user-id",
				Issuer:    "liora-refresh",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testSecret))
		require.NoError(t, err)

		wrongSecret := "wrong-secret-key"
		validatedClaims, err := server.validateRefreshToken(tokenString, wrongSecret)

		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("Expired Token", func(t *testing.T) {
		claims := &types.CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "test-user-id",
				Issuer:    "liora-refresh",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 1 hour ago
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testSecret))
		require.NoError(t, err)

		validatedClaims, err := server.validateRefreshToken(tokenString, testSecret)

		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
		assert.Contains(t, err.Error(), "token expired")
	})

	t.Run("Malformed Token", func(t *testing.T) {
		malformedToken := "not.a.valid.jwt.token"

		validatedClaims, err := server.validateRefreshToken(malformedToken, testSecret)

		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("Empty Token", func(t *testing.T) {
		validatedClaims, err := server.validateRefreshToken("", testSecret)

		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("Wrong Signing Method", func(t *testing.T) {
		claims := &types.CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "test-user-id",
				Issuer:    "liora-refresh",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		fmt.Println(token.Header)

		tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature"

		validatedClaims, err := server.validateRefreshToken(tokenString, testSecret)

		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("Token Valid But Claims Invalid", func(t *testing.T) {
		standardClaims := jwt.RegisteredClaims{
			Subject:   "test-user-id",
			Issuer:    "liora-refresh",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, standardClaims)
		tokenString, err := token.SignedString([]byte(testSecret))
		require.NoError(t, err)

		validatedClaims, err := server.validateRefreshToken(tokenString, testSecret)

		assert.NoError(t, err)
		assert.NotNil(t, validatedClaims)
		assert.Equal(t, "test-user-id", validatedClaims.Subject)
	})
}

func BenchmarkValidateRefreshToken(b *testing.B) {
	testSecret := "test-secret-key-for-jwt-signing"
	server := &Server{}

	claims := &types.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "test-user-id",
			Issuer:    "liora-refresh",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testSecret))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = server.validateRefreshToken(tokenString, testSecret)
	}
}
