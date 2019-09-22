package graphql

import (
	"context"
	"net/http"
)

func Authorization() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("authorization")
			ctx := context.WithValue(r.Context(), "authToken", auth)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
