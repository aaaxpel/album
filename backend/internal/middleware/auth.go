package middleware

import (
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apparently need to define your own type to avoid collisions.
		// Will check this later when adding auth endpoints.
		// ctx := context.WithValue(r.Context(), "user", "test")

		// next.ServeHTTP(w, r.WithContext(ctx))
	})
}
