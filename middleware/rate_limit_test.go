package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCache struct {
	data map[string]int
}

func (m *mockCache) GetIncrement(key string) (int, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}

	return 1, nil
}

func (m *mockCache) SetExpire(key string, value int, duration int) error {
	m.data[key] = value
	return nil
}

func (m *mockCache) Increment(key string) error {
	m.data[key]++
	return nil
}

func TestRateLimiterMiddleware(t *testing.T) {
	os.Setenv("TOKEN_MAX_REQUESTS", "5")
	os.Setenv("TOKEN_DURATION", "60")
	os.Setenv("IP_MAX_REQUESTS", "5")
	os.Setenv("IP_DURATION", "60")

	cache := &mockCache{data: make(map[string]int)}
	handler := RateLimiterMiddleware(cache, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("allows requests under the limit", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "test_token")
		rr := httptest.NewRecorder()

		for i := 0; i < 5; i++ {
			handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
		}
	})

	t.Run("blocks requests over the limit", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "test_token")
		rr := httptest.NewRecorder()

		for i := 0; i < 6; i++ {
			rr = httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
		}

		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	})

	t.Run("allows requests under the limit by IP", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.RemoteAddr = "127.0.0.1"

		for i := 0; i < 5; i++ {
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
		}
	})

	t.Run("blocks requests over the limit by IP", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.RemoteAddr = "127.0.0.1"
		rr := httptest.NewRecorder()

		for i := 0; i < 6; i++ {
			rr = httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
		}
		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	})
}
