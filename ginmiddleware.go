package ginratelimit

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RateLimitByIP(tb *TokenBucket) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		rateLimit(ctx, tb, ip)
	}
}

func RateLimitByUserId(tb *TokenBucket, userId string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rateLimit(ctx, tb, userId)
	}
}

func PreventBruteForce(tb *TokenBucket, userKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		rateLimit(ctx, tb, ip)
		rateLimit(ctx, tb, userKey)
	}
}

func rateLimit(ctx *gin.Context, tb *TokenBucket, key string) {
	if !tb.Allow(key) {
		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		return
	}
	ctx.Next()
}
