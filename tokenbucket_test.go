package ginratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewTokenBucket(t *testing.T) {
	tb := NewTokenBucket(10, 1*time.Minute)
	if tb.threshold != 10 {
		t.Errorf("Expected threshold to be 10, got %d", tb.threshold)
	}
	if tb.ttl != 1*time.Minute {
		t.Errorf("Expected ttl to be 1 minute, got %s", tb.ttl)
	}
	if len(tb.tokens) != 0 {
		t.Errorf("Expected tokens to be empty, got %d", len(tb.tokens))
	}
}

func TestTokenBucket_Allow(t *testing.T) {
	tb := NewTokenBucket(2, 50*time.Millisecond)

	// Test for a single key
	key := "testKey"

	if !tb.Allow(key) {
		t.Errorf("Expected to allow the first request for key %s", key)
	}

	if !tb.Allow(key) {
		t.Errorf("Expected to allow the second request for key %s", key)
	}

	if tb.Allow(key) {
		t.Errorf("Expected to reject the third request for key %s", key)
	}

	// Wait for tokens to refill
	time.Sleep(60 * time.Millisecond)

	if !tb.Allow(key) {
		t.Errorf("Expected to allow request after ttl for key %s", key)
	}
}

func TestTokenBucket_Concurrency(t *testing.T) {
	tb := NewTokenBucket(10, 10*time.Millisecond)
	key := "concurrencyKey"

	var wg sync.WaitGroup
	var allowedRequests int32 = 0
	totalRequests := 20

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if tb.Allow(key) {
				atomic.AddInt32(&allowedRequests, 1)
			}
		}()
	}

	wg.Wait()

	allowedRequestsFinal := atomic.LoadInt32(&allowedRequests)
	if allowedRequestsFinal > int32(tb.threshold) {
		t.Errorf("Allowed more requests (%d) than threshold (%d)", allowedRequestsFinal, tb.threshold)
	}

	// Wait for tokens to refill
	time.Sleep(20 * time.Millisecond)

	// Reset for after refill
	atomic.StoreInt32(&allowedRequests, 0)

	for i := 0; i < totalRequests; i++ {
		if tb.Allow(key) {
			atomic.AddInt32(&allowedRequests, 1)
		}
	}

	allowedAfterRefill := atomic.LoadInt32(&allowedRequests)
	if allowedAfterRefill != int32(tb.threshold) {
		t.Errorf("Expected to allow %d requests after refill, but allowed %d", tb.threshold, allowedAfterRefill)
	}
}
