package middleware

import (
	"backend/cmd/helper"
	"net/http"
	"sync"
	"log/slog"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	Limiters map[string]*rate.Limiter
	Mu sync.Mutex
	Rps rate.Limit
	Burst int
}

func NewIPRateLimiter(rps rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		Limiters: make(map[string]*rate.Limiter),
		Rps: rps,
		Burst: burst,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.Mu.Lock()
	defer i.Mu.Unlock()

	limiter, exists := i.Limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(i.Rps, i.Burst)
		i.Limiters[ip] = limiter
	}

	return limiter
}

func (i *IPRateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := helper.GetClientIP(r)
		limiter := i.GetLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		slog.Info("[REQUEST]", "method", r.Method, "path", r.URL.Path, "ip", ip)
		next.ServeHTTP(w, r)
	})
}