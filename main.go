package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/wendellnd/go-rate-limiter-challenge/cache"
	"github.com/wendellnd/go-rate-limiter-challenge/middleware"
)

var (
	redisClient cache.Cache
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	redisClient = cache.NewRedis()
}

func NewServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: middleware.RateLimiterMiddleware(redisClient, mux),
	}
}

func main() {
	server := NewServer()

	log.Println("Server running on port 8080")
	log.Fatal(server.ListenAndServe())
}
