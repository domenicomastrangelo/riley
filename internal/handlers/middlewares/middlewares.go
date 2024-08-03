package middlewares

import "net/http"

func DefaultMiddlewares(handler http.HandlerFunc) http.Handler {
	return Authorization(
		Authentication(
			RateLimiter(
				handler,
			),
		),
	)
}
