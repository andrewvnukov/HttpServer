package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Middleware для проверки API ключа
func APIKeyMiddleware(validAPIKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем API ключ из заголовка
			apiKey := r.Header.Get("X-API-Key")

			// Если нет в заголовке, проверяем query параметр
			if apiKey == "" {
				apiKey = r.URL.Query().Get("api_key")
			}

			// Проверяем ключ
			if apiKey != validAPIKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Invalid or missing API key"}`))
				return
			}
			// Ключ верный, продолжаем обработку
			next.ServeHTTP(w, r)
		})
	}
}
