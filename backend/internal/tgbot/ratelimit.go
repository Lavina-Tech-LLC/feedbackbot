package tgbot

import (
	"sync"
	"time"
)

const maxFeedbackPerHour = 10

type rateLimiter struct {
	mu      sync.Mutex
	entries map[int64][]time.Time
}

var limiter = &rateLimiter{
	entries: make(map[int64][]time.Time),
}

// checkRateLimit returns true if the user is within limits
func checkRateLimit(userID int64) bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-1 * time.Hour)

	// Clean old entries
	timestamps := limiter.entries[userID]
	var valid []time.Time
	for _, t := range timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= maxFeedbackPerHour {
		limiter.entries[userID] = valid
		return false
	}

	limiter.entries[userID] = append(valid, now)
	return true
}
