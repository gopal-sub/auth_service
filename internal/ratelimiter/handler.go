package ratelimiter


import (
	"net/http"
)

func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Logic executed before the core handler runs (e.g., Auth, Logging)
        w.Header().Set("X-Custom-Header", "Golang")

        // 2. Trigger the next step in the request chain
        next.ServeHTTP(w, r)

        // 3. Logic executed after the core handler completes
    })
}