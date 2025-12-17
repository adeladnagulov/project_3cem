package middleware

import "net/http"

var address string = "http://5b7967c6-f561-4a5d-80f1-64bcfef366e1.tunnel4.com"

func CORSmiddlewera(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", address)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		/*
			if r.Method == "OPTIONS" {
				// Разрешаем необходимые методы и заголовки
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Subdomain, Authorization")
				w.WriteHeader(http.StatusOK)
				return // Важно: завершаем обработку для OPTIONS
			}
		*/
		next.ServeHTTP(w, r)
	})
}
