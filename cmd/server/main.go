package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/JeanCarlos20-code/rateLimit/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	tokens   = make(map[string]*rate.Limiter) // Limite por token
	mu       sync.Mutex
)

func getIp(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println("Error getting IP:", err)
		return ""
	}
	return ip
}

func getVisitorByIP(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(1, 1)
		visitors[ip] = limiter
	}
	return limiter
}

func getVisitorByToken(token string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := tokens[token]
	if !exists {
		limiter = rate.NewLimiter(10, 1)
		tokens[token] = limiter
	}
	return limiter
}

func rateLimitByIPAndToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIp(r)
		token := r.Header.Get("API_KEY")

		ipLimiter := getVisitorByIP(ip)
		tokenLimiter := getVisitorByToken(token)

		if !ipLimiter.Allow() || !tokenLimiter.Allow() {
			log.Printf("[BLOCKED] Request limit exceeded from IP %s with token %s", ip, token)
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		log.Printf("[ALLOWED] Request from IP %s with token %s", ip, token)
		h.ServeHTTP(w, r)
	})
}

// func requestIpDifferent() {
// 	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
// 	for _, ip := range ips {
// 		wg.Add(1)
// 		go func(ip string) {
// 			defer wg.Done()

// 			req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
// 			req.RemoteAddr = ip + ":12345"
// 			_, err := http.DefaultClient.Do(req)
// 			if err != nil {
// 				log.Println("Error in request:", err)
// 			}
// 		}(ip)
// 	}
// 	wg.Wait()
// }

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.GetTokenAuth()))
	r.Use(middleware.WithValue("JwtExpiresIn", config.GetJWTExpiresIn()))

	r.Use(jwtauth.Verifier(config.GetTokenAuth()))
	r.Use(jwtauth.Authenticator)

	r.Get("/teste", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	http.ListenAndServe(":8080", rateLimitByIPAndToken(r))
}
