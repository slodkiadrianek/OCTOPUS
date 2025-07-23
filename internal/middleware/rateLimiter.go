package middleware

import "time"

type RateLimiter struct {
	Token      int
	MaxTokens  int
	refillRate time.Duration
	lastRefill time.Duration
}
