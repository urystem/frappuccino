package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func Middleware(handler http.HandlerFunc) http.HandlerFunc {
	handler = WrapContext(handler)

	return handler
}

func WrapContext(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
			defer slog.Debug("query took more than 5 seconds", r)
		})
}
