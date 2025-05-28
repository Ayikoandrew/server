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
			HttpOnly: true,
			Secure:   true,
			MaxAge:   30 * 60,
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
			Secure:   false,
			MaxAge:   7 * 24 * 60 * 60,
			SameSite: http.SameSiteStrictMode,
		},
	)
}

func ClearTokenCookies(w http.ResponseWriter) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			MaxAge:   -1,
			SameSite: http.SameSiteStrictMode,
		},
	)
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			MaxAge:   -1,
			SameSite: http.SameSiteStrictMode,
		},
	)
}
