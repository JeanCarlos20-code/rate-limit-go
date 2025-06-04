package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestServer(ipLimit, tokenLimit, blockSeconds int) http.Handler {
	limiter := NewLimiter(ipLimit, tokenLimit, blockSeconds)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return Middleware(limiter)(mux)
}

func TestMiddleware_IPBlocking(t *testing.T) {
	server := setupTestServer(1, 100, 1)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("1st IP request expected 200, got %d", rec.Code)
	}

	// 2ª deve ser bloqueada
	rec2 := httptest.NewRecorder()
	server.ServeHTTP(rec2, req)

	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("2nd IP request expected 429, got %d", rec2.Code)
	}
}

func TestMiddleware_TokenPriority(t *testing.T) {
	server := setupTestServer(1, 2, 1)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "4.3.2.1:8888"
	req.Header.Set("API_KEY", "abc123")

	rec1 := httptest.NewRecorder()
	server.ServeHTTP(rec1, req)
	if rec1.Code != http.StatusOK {
		t.Errorf("Token 1st request expected 200, got %d", rec1.Code)
	}

	rec2 := httptest.NewRecorder()
	server.ServeHTTP(rec2, req)
	if rec2.Code != http.StatusOK {
		t.Errorf("Token 2nd request expected 200, got %d", rec2.Code)
	}

	rec3 := httptest.NewRecorder()
	server.ServeHTTP(rec3, req)
	if rec3.Code != http.StatusTooManyRequests {
		t.Errorf("Token 3rd request expected 429, got %d", rec3.Code)
	}
}

func TestMiddleware_BlockExpires(t *testing.T) {
	server := setupTestServer(1, 100, 1)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "6.6.6.6:666"

	rec1 := httptest.NewRecorder()
	server.ServeHTTP(rec1, req)
	if rec1.Code != http.StatusOK {
		t.Error("1st request expected 200")
	}

	rec2 := httptest.NewRecorder()
	server.ServeHTTP(rec2, req)
	if rec2.Code != http.StatusTooManyRequests {
		t.Error("2nd request expected 429 (blocked)")
	}

	time.Sleep(1100 * time.Millisecond)

	rec3 := httptest.NewRecorder()
	server.ServeHTTP(rec3, req)
	if rec3.Code != http.StatusOK {
		t.Error("3rd request (after expiration) expected 200")
	}
}
