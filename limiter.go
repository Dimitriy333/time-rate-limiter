package time_rate_limiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	counts   map[string]int
	timeouts map[string]time.Time
	done     chan struct{}
}

func NewRateLimiter(cleanUpInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		counts:   make(map[string]int),
		timeouts: make(map[string]time.Time),
		done:     make(chan struct{}),
	}
	go rl.cleanUp(cleanUpInterval) // start clean up routine
	return rl
}

func (rl *RateLimiter) Limit(key string, limit int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// check timeout for given key
	if rl.timeouts[key].Before(now) {
		rl.counts[key] = 0
		rl.timeouts[key] = now.Add(window)
	}

	// check limit for given key
	if rl.counts[key] >= limit {
		return false
	}

	// increment counter if request passed
	rl.counts[key]++
	return true
}

func (rl *RateLimiter) cleanUp(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, expiry := range rl.timeouts {
				if expiry.Before(now) {
					// remove expired keys
					delete(rl.counts, key)
					delete(rl.timeouts, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.done: // finish cycle when done signal received
			return
		}
	}
}

func (rl *RateLimiter) Dispose() {
	close(rl.done) // finish clean up routine
}
