package middleware

import (
	"net/http"
	"time"
)

type RateLimiter struct {
	Tokens        int
	MaxTokens     int
	RefillRate    time.Duration
	LastRefill    time.Time
	BlockUntil    time.Time
	BlockDuration time.Duration
}

func RateLimiterHandler (next http.Handler) http.Handler{
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    ipAddress:= r.RemoteAddr()
    if  
    next.ServeHTTP(w,r)
  })
}

func NewRateLimiter(maxTokens int, refillRate time.Duration, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		Tokens:        maxTokens,
		MaxTokens:     maxTokens,
		RefillRate:    refillRate,
		LastRefill:    time.Now(),
		BlockDuration: blockDuration,
	}
}


func (ra *RateLimiter) Allow() bool {
	now := time.Now()
	if now.Before(ra.BlockUntil) {
		return false
	}
	elapsed := now.Sub(ra.LastRefill)
	newTokens := (elapsed / ra.RefillRate)
	if newTokens > 0 {
		ra.Tokens += int(newTokens)
		if ra.Tokens > ra.MaxTokens {
			ra.Tokens = ra.MaxTokens
		}
		ra.LastRefill = now
	}

	if ra.Tokens > 0 {
		ra.Tokens--
		return true
	}

	return false
}
