package middleware

import (
	"cafeteria/pkg/jtoken"
	"context"
	"net/http"
)

type key string

const ctxKey key = "payload"

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwtToken")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Not authorized", http.StatusForbidden)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		jwtToken := cookie.Value
		payload, err := jtoken.DecodePayload(jwtToken)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKey, payload)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}
