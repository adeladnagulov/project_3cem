package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
)

const (
	SubdomainKey ctxKey = "subdomain"
)

func SubdomainMiddlewera(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		log.Printf("host: %s", host)

		host = strings.Split(host, ":")[0]
		parts := strings.Split(host, ".")
		if len(parts) >= 2 {
			sub := parts[0]
			if sub != "localhost" && sub != "127" {
				ctx := context.WithValue(r.Context(), SubdomainKey, sub)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
