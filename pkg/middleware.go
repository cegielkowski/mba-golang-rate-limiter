package pkg

import (
	"context"
	"net/http"
)

func RateLimitMiddleware(rl *RateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			token := r.Header.Get("API_KEY")

			var allowed bool
			var err error

			ctx := context.Background()

			if token != "" {
				allowed, err = rl.AllowToken(ctx, token)
			} else {
				allowed, err = rl.AllowIP(ctx, ip)
			}

			if err != nil || !allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
