package limiter

import (
	"testing"
	"time"
)

func TestLimiter_AllowsRequests_UnderIPLimit(t *testing.T) {
	lim := NewLimiter(3, 10, 1) // IP: 3 reqs, Token: 10 reqs, 1s bloqueio
	ip := "192.168.1.1"

	for i := range 3 {
		if !lim.Allow(ip, "") {
			t.Errorf("Request %d from IP %s should have been allowed", i+1, ip)
		}
	}
}

func TestLimiter_BlocksRequests_OverIPLimit(t *testing.T) {
	lim := NewLimiter(2, 10, 2)
	ip := "10.0.0.1"

	lim.Allow(ip, "")
	lim.Allow(ip, "")
	allowed := lim.Allow(ip, "") // terceira, deve bloquear

	if allowed {
		t.Error("Expected request to be blocked after exceeding IP limit")
	}
}

func TestLimiter_AllowsRequests_AfterBlockExpires(t *testing.T) {
	lim := NewLimiter(1, 10, 1) // 1 por IP, 1s de bloqueio
	ip := "127.0.0.1"

	_ = lim.Allow(ip, "") // ok
	_ = lim.Allow(ip, "") // bloqueia

	time.Sleep(1100 * time.Millisecond)

	if !lim.Allow(ip, "") {
		t.Error("Expected request to be allowed after block duration expired")
	}
}

func TestLimiter_UsesTokenLimit_WhenProvided(t *testing.T) {
	lim := NewLimiter(1, 3, 2)
	ip := "127.0.0.2"
	token := "my-token"

	for i := range 3 {
		if !lim.Allow(ip, token) {
			t.Errorf("Token request %d should be allowed", i+1)
		}
	}

	if lim.Allow(ip, token) {
		t.Error("Token request 4 should be blocked")
	}
}

func TestLimiter_TokenLimit_DoesNotInterfereWithIP(t *testing.T) {
	lim := NewLimiter(1, 2, 1)
	ip := "192.168.0.1"

	// limite IP (1 req)
	if !lim.Allow(ip, "") {
		t.Fatal("1ª requisição IP deveria passar")
	}
	if lim.Allow(ip, "") {
		t.Fatal("2ª requisição IP deveria ser bloqueada")
	}

	// Token novo
	if !lim.Allow(ip, "token-x") {
		t.Fatal("Token-x 1ª requisição deveria ser permitida")
	}
	if !lim.Allow(ip, "token-x") {
		t.Fatal("Token-x 2ª requisição deveria ser permitida")
	}
	if lim.Allow(ip, "token-x") {
		t.Fatal("Token-x 3ª requisição deveria ser bloqueada")
	}
}

func TestLimiter_Stress(t *testing.T) {
	lim := NewLimiter(5, 10, 2)

	ip := "192.168.0.100"
	token := "stress-token"

	successIP := 0
	for range 10 {
		if lim.Allow(ip, "") {
			successIP++
		}
	}
	if successIP != 5 {
		t.Errorf("Esperava 5 requisições IP permitidas, teve %d", successIP)
	}

	successToken := 0
	for range 15 {
		if lim.Allow("0.0.0.0", token) {
			successToken++
		}
	}
	if successToken != 10 {
		t.Errorf("Esperava 10 requisições Token permitidas, teve %d", successToken)
	}

	time.Sleep(2100 * time.Millisecond) // Desbloqueia o token para ser testado

	successBoth := 0
	for range 15 {
		if lim.Allow(ip, token) {
			successBoth++
		}
	}
	if successBoth != 10 {
		t.Errorf("Esperava 10 requisições IP+Token permitidas (via token), teve %d", successBoth)
	}
}
