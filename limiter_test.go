package time_rate_limiter

import (
	"testing"
	"time"
)

func TestRateLimiter_Limit(t *testing.T) {
	rl := NewRateLimiter(time.Minute)
	dur := 100 * time.Millisecond

	tests := []struct {
		name         string
		key          string
		limit        int
		window       time.Duration
		prevDuration time.Duration
		count        int
		newCount     int
		want         bool
	}{
		{"BelowLimit", "1", 5, dur, dur, 4, 5, true},
		{"AtLimit", "2", 5, dur, dur, 5, 5, false},
		{"AboveLimit", "3", 5, dur, dur, 6, 6, false},
		{"DifferentKeys", "4", 3, dur, dur, 2, 3, true},
		{"DifferentWindows", "5", 5, 2 * dur, dur, 4, 5, true},
		{"ExpiredTimeout", "6", 5, dur, -2 * dur, 3, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set initial counts and timeouts
			rl.mu.Lock()
			rl.counts[tt.key] = tt.count
			rl.timeouts[tt.key] = time.Now().Add(tt.prevDuration)
			rl.mu.Unlock()

			// Perform the test
			if got := rl.Limit(tt.key, tt.limit, tt.window); got != tt.want {
				t.Errorf("Limit() for key %s, count %d, window %s = %v; want %v",
					tt.key, tt.count, tt.window, got, tt.want)
			}

			// If the expected result is true, check if the count has increased
			if tt.want && rl.counts[tt.key] != tt.newCount {
				t.Errorf("Limit() for key %s, count %d, window %s did not increase count",
					tt.key, tt.count, tt.window)
			}
		})
	}
}

func BenchmarkRateLimiter_Limit(b *testing.B) {
	rl := NewRateLimiter(time.Minute)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rl.Limit("user", 5, time.Second)
	}
}
