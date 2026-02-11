package middleware

import (
	"net"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// rateLimiterStore manages rate limiters for different IP addresses.
type rateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
}

func NewRateLimiterStore(r rate.Limit, burst int) *rateLimiterStore {
	return &rateLimiterStore{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
	}
}

// getLimiter retrieves the rate limiter for the given IP address, creating a new one if it doesn't exist.
func (s *rateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, exists := s.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(s.rate, s.burst)
		s.limiters[ip] = limiter
	}

	return limiter
}

// Middleware is an HTTP middleware that applies rate limiting based on the client's IP address.
func (s *rateLimiterStore) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "unable to determine client IP", http.StatusInternalServerError)
			return
		}

		limiter := s.getLimiter(ip)

		if !limiter.Allow() {
			logrus.WithFields(logrus.Fields{
				"component": "middleware",
				"middleware": "rate_limiter",
				"client_ip": ip,
				"path":      r.URL.Path,
				"method":    r.Method,
			}).Warn("rate limit exceeded")
			w.Header().Set("Retry-After", "60")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
