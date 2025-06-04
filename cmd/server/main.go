package main

import (
	"log"
	"net/http"

	"github.com/JeanCarlos20-code/rateLimit/config"
	"github.com/JeanCarlos20-code/rateLimit/internal/limiter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	rl := limiter.NewLimiter(cfg.GetIPLimit(), cfg.GetTokenLimit(), cfg.GetBlockDuration())

	if cfg.GetRateBackend() == "redis" {
		log.Println("Rate limiter using Redis")
		rl.UseRedis(cfg.GetRedisAddr(), cfg.GetRedisPassword())
	} else {
		log.Println("Rate limiter using local in-memory")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(limiter.Middleware(rl))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello! You have successfully accessed within the allowed limit."))
	})

	log.Println("Server is running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
