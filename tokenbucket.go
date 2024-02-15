package ginratelimit

import (
	"sync"
	"time"
)

type TokenInfo struct {
	RemainingTokens int
	ExpiresAt       time.Time
}

type TokenBucket struct {
	threshold int
	ttl       time.Duration
	tokens    map[string]TokenInfo
	mu        sync.Mutex
}

func NewTokenBucket(threshold int, ttl time.Duration) *TokenBucket {
	return &TokenBucket{
		threshold: threshold,
		ttl:       ttl,
		tokens:    make(map[string]TokenInfo),
	}
}

func (tb *TokenBucket) Allow(key string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	tokenInfo, exists := tb.tokens[key]

	// Refill tokens if the record exists and has expired
	if exists && now.After(tokenInfo.ExpiresAt) {
		tokenInfo = TokenInfo{
			RemainingTokens: tb.threshold,
			ExpiresAt:       now.Add(tb.ttl),
		}
	}

	// If the key doesn't exist, initialize it
	if !exists {
		tokenInfo = TokenInfo{
			RemainingTokens: tb.threshold - 1, // Consume a token
			ExpiresAt:       now.Add(tb.ttl),
		}
		tb.tokens[key] = tokenInfo
		return true
	}

	// If tokens are available, decrement and allow
	if tokenInfo.RemainingTokens > 0 {
		tokenInfo.RemainingTokens--
		tb.tokens[key] = tokenInfo
		return true
	}

	// No tokens available, reject
	return false
}
