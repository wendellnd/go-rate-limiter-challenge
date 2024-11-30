package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/wendellnd/go-rate-limiter-challenge/cache"
)

func handleRateLimit(cache cache.Cache, maxReq int, duration int, key string) error {
	key = fmt.Sprintf("rate_limit:%s", key)
	count, err := cache.GetIncrement(key)
	if err != nil {
		return err
	}

	if count > maxReq {
		return errors.New("you have reached the maximum number of requests or actions allowed within a certain time frame")
	}

	cache.Increment(key)
	cache.SetExpire(key, count+1, duration)

	return nil
}

func RateLimiterMiddleware(cache cache.Cache, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")

		var err error
		if token != "" {
			maxReq := os.Getenv("TOKEN_MAX_REQUESTS")
			duration := os.Getenv("TOKEN_DURATION")
			maxReqInt, _ := strconv.Atoi(maxReq)
			durationInt, _ := strconv.Atoi(duration)

			err = handleRateLimit(cache, maxReqInt, durationInt, token)
		} else {
			ip := r.RemoteAddr

			maxReq := os.Getenv("IP_MAX_REQUESTS")
			duration := os.Getenv("IP_DURATION")
			maxReqInt, _ := strconv.Atoi(maxReq)
			durationInt, _ := strconv.Atoi(duration)

			err = handleRateLimit(cache, maxReqInt, durationInt, ip)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
