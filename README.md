# Token Bucket Rate Limiter for Gin-Gonic

This library implements a Token Bucket Algorithm-based Rate Limiter designed for easy integration with the Gin-Gonic web framework. It provides an efficient way to control the rate of incoming requests, ensuring your application remains reliable and responsive under varying loads.

## Features

- **Token Bucket Rate Limiting**: Utilizes the token bucket algorithm to allow for a configurable rate of requests, accommodating bursts of traffic efficiently.
- **Seamless Integration with Gin-Gonic**: Designed as middleware that can be easily added to your Gin-Gonic routes or router groups.
- **Customizable Limits**: Offers the flexibility to set request thresholds and token refill intervals per route, adapting to your application's specific needs.

## Installation

To install the rate limiter middleware, use the `go get` command:

```sh
go get github.com/ljahier/gin-ratelimit
```

## Usage

Import the rate limiter package into your Gin application, and then initialize and apply the middleware as shown in the example below.

### Basic Setup

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ljahier/gin-ratelimit"
    "time"
)

func main() {
    // Initialize the Gin router
    r := gin.Default()

    // Create a new token bucket rate limiter
    tb := ginratelimit.NewTokenBucket(100, 1*time.Minute) // 100 requests per minute

    // Apply the rate limiter middleware to all routes
    r.Use(ginratelimit.RateLimitByIP(tb))

    // Define a example route
    r.GET("/example", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Rate limited request succeeded!",
        })
    })

    // Start the Gin server
    r.Run() // Listen and serve on 0.0.0.0:8080
}
```

### Rate Limiting by User id

If you need to rate limit based on a user id (e.g., extracted from a JWT token or an API key), you can use the `RateLimitByUserId` middleware. Here's how to set it up:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ljahier/go-gin-ratelimit"
	"time"
)

// Assuming you have a function that authenticates the user
func Authenticate(ctx *gin.Context) {
    // ... your authenticate logic
    ctx.Set("userId", "xxx-yyy-zzz")
    ctx.Next()
}

// Assuming you have a function to extract the user id
func extractUserId(ctx *gin.Context) string {
    // Extract the user id from the request, e.g., from headers or JWT token
    return ctx.GetString("userId")
}

func main() {
    r := gin.Default()
    
    // Initialize the token bucket rate limiter
    tb := ginratelimit.NewTokenBucket(50, 1*time.Minute) // 50 requests per minute per user
    
    r.Use(Authenticate)
    
    // Apply the rate limiter middleware using a custom user id extractor
    r.Use(func(ctx *gin.Context) {
        userId := extractUserId(ctx)
        ginratelimit.RateLimitByUserId(tb, userId)(ctx)
    })
    
    r.GET("/user-specific-route", func(c *gin.Context) {
    c.JSON(200, gin.H{
            "message": "User-specific rate limited request succeeded!",
        })
    })
    
    r.Run(":9090")
}

```

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues to improve the library or fix problems.

## License

This project is licensed under the MIT License - see the LICENSE file for details.