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
		log.Fatal(err)
	}
	rateLimiter := limiter.NewLimiter(5, cfg.GetJWTExpiresIn(), 60)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(limiter.Middleware(rateLimiter))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome!"))
	})

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
