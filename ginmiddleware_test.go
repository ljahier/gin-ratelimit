package ginratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// Function to create a ping router with middleware
func setupRouter(tb *TokenBucket, middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return r
}

func TestRateLimitByIP(t *testing.T) {
	tb := NewTokenBucket(1, 1*time.Minute) // Allow 1 request per minute per IP
	r := setupRouter(tb, RateLimitByIP(tb))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.RemoteAddr = "127.0.0.1:12345" // Set client IP

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Second request should be rate limited
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}
}

func TestRateLimitByUserId(t *testing.T) {
	tb := NewTokenBucket(1, 1*time.Minute) // Allow 1 request per minute per userId
	userId := "testUser"
	r := setupRouter(tb, RateLimitByUserId(tb, userId))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Second request should be rate limited
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}
}

func TestPreventBruteForce(t *testing.T) {
	tb := NewTokenBucket(1, 1*time.Minute) // Allow 1 request per minute per IP and per userId
	userId := "testUser"
	r := setupRouter(tb, PreventBruteForce(tb, userId))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.RemoteAddr = "127.0.0.1:12345" // Set client IP

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Second request should be rate limited by IP
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}

	// Reset recorder and request for user Id rate limit
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 due to userId limit, got %d", w.Code)
	}
}

func TestPreventBruteForce2(t *testing.T) {
	tb := NewTokenBucket(1, 1*time.Minute) // Allow 1 request per minute per IP and per userId
	userEmail := "john@doe.com"
	r := setupRouter(tb, PreventBruteForce(tb, userEmail))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.RemoteAddr = "127.0.0.1:12345" // Set client IP

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Second request should be rate limited by IP
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}

	// Reset recorder and request for user Id rate limit
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 due to userEmail limit, got %d", w.Code)
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	// Setup a token bucket with a small TTL for testing
	tb := NewTokenBucket(1, 2*time.Second) // 1 request allowed, with a 2-second TTL
	r := setupRouter(tb, RateLimitByIP(tb))

	// First request should be allowed
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.RemoteAddr = "192.168.1.1:12345" // Set client IP
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("First request: Expected status 200, got %d", w.Code)
	}

	// Immediately try a second request, which should be rate-limited
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Second request: Expected status 429, got %d", w.Code)
	}

	// Wait for longer than the TTL for the bucket to refill
	time.Sleep(3 * time.Second)

	// Try another request, which should now be allowed
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Third request after TTL: Expected status 200, got %d", w.Code)
	}
}
