package middlewares

import (
	"net/http"

	"riley/internal/auth"
	"riley/internal/config"
)

func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		secret := config.LoadConfig().TokenSecret

		userID, err := auth.GetUserIDFromToken(token, secret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check that user is authorized to access the resource
		if !auth.IsAuthorized(r.URL, userID) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next(w, r)
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
