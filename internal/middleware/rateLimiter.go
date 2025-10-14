package middleware

import (
	"fmt"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"net/http"
	"time"
)

type RateLimiterUser struct {
	Tokens     int
	LastRefill time.Time
	BlockUntil time.Time
}

type RateLimiter struct {
	Users           map[string]*RateLimiterUser
	LoggerService *utils.Logger
	MaxTokens       int
	RefillRate      time.Duration
	BlockDuration   time.Duration
	InactiveTimeout time.Duration
}

func RateLimiterMiddleware(userMap map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RateLimiterHandler(next, userMap)
	}
}

func RateLimiterHandler(next http.Handler, rateLimiter RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.RemoteAddr
		if _, exists := rateLimiter.Users[ipAddress]; exists {
			res := rateLimiter.Allow(ipAddress)
			if res {
				next.ServeHTTP(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
		return
	})
}

func NewRateLimiter(maxTokens int, refillRate time.Duration, blockDuration time.Duration, inactiveTimeout time.Duration, loggerService *utils.Logger) *RateLimiter {
	return &RateLimiter{
		Users:           map[string]*RateLimiterUser{},
		LoggerService: loggerService,
		MaxTokens:       maxTokens,
		RefillRate:      refillRate,
		BlockDuration:   blockDuration,
		InactiveTimeout: inactiveTimeout,
	}
}

func (ra *RateLimiter) Cleanup() {
	now := time.Now()
	removeTime := now.Add(ra.InactiveTimeout)
	removed := 0
	for ipAddress, user := range ra.Users {
		if removeTime.Before(user.LastRefill) {
			delete(ra.Users, ipAddress)
			removed++
		}
	}
	ra.LoggerService.Info("Removed users from rate limiter",removed)
}

func (ra *RateLimiter) Allow(ipAddress string) bool {
	now := time.Now()
	user, _ := ra.Users[ipAddress]
	if now.Before(user.BlockUntil) {
		return false
	}
	elapsed := now.Sub(user.LastRefill)
	newTokens := elapsed / ra.RefillRate
	if newTokens > 0 {
		user.Tokens += int(newTokens)
		if user.Tokens > ra.MaxTokens {
			user.Tokens = ra.MaxTokens
		}
		user.LastRefill = now
	}

	if user.Tokens > 0 {
		user.Tokens--
		return true
	}

	return false
}
