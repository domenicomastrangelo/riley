package middlewares

import (
	"net/http"
	"riley/internal/auth"
	"riley/internal/config"
)

func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func Authentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		secret := config.LoadConfig().TokenSecret

		if err := auth.CheckToken(token, secret); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}
