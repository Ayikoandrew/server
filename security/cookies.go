package security

import (
	"net/http"
	"os"
)

func SetTokenCookies(w http.ResponseWriter, accessToken, refreshToken string) {

	isProduction := os.Getenv("ENV") == "production"
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Path:     "/",
			HttpOnly: false,
			Secure:   isProduction,
			MaxAge:   15 * 60,
			SameSite: http.SameSiteStrictMode,
		},
	)

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   isProduction,
			MaxAge:   7 * 24 * 60 * 60,
			SameSite: http.SameSiteStrictMode,
		},
	)
}
