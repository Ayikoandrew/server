package security

import (
	"net/http"
)

func SetTokenCookies(w http.ResponseWriter, accessToken, refreshToken string) {

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Path:     "/",
			HttpOnly: false,
			Secure:   true,
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
			Secure:   true,
			MaxAge:   7 * 24 * 60 * 60,
			SameSite: http.SameSiteStrictMode,
		},
	)
}
