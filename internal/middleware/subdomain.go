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
		subdomain := r.Header.Get("X-Subdomain")

		if subdomain == "" {
			host := r.Header.Get("X-Forwarded-Host")
			if host == "" {
				host = r.Host
			}

			if host != "" {
				host = strings.Split(host, ":")[0]
				parts := strings.Split(host, ".")

				if len(parts) >= 2 && parts[1] == "localhost" {
					subdomain = parts[0]
				} else if len(parts) >= 3 && strings.Contains(host, "tunnel4.com") {
					subdomain = parts[0]
				}
			}
		}

		log.Printf("Subdomain: '%s'", subdomain)

		ctx := context.WithValue(r.Context(), SubdomainKey, subdomain)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
