# Time Rate Limiter

This is a simple time rate limiter implemented in Go. It allows you to restrict the rate of certain actions based on keys (e.g., user IDs, IP addresses) within a specified time window.

### Features
- Allows you to set a rate limit for specific keys within a time window.
- Thread-safe implementation using mutex locks to handle concurrent access.
- Provides a Dispose method to stop the internal goroutine for cleanup.

## Installation

You can include this rate limiter in your Go project by importing the package:

```bash
go get github.com/Dimitriy333/time-rate-limiter
```

## Usage
```
import (
    rate_limiter "github.com/Dimitriy333/time-rate-limiter"
)

// Create a new rate limiter
limiterCleanUpInterval := 5 * time.Minute
rl := rate_limiter.NewRateLimiter(limiterCleanUpInterval)

// Check if the rate limit for a key is exceeded
key := "user123"
limit := 5
window := time.Minute
if rl.Limit(key, limit, window) {
    // Allow the action
} else {
    // Deny the action
}
```
