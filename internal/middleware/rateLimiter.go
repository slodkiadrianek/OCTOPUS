package middleware

import (
	"context"
	"net/http"
	"strings"

	// "strings"
	"time"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type RateLimiterUser struct {
	Tokens     int
	LastRefill time.Time
	BlockUntil time.Time
}

func NewRateLimiterUser(token int, lastRefill time.Time, blockUntil time.Time) *RateLimiterUser {
	return &RateLimiterUser{
		Tokens:     token,
		LastRefill: lastRefill,
		BlockUntil: blockUntil,
	}
}

type RateLimiter struct {
	Users           map[string]*RateLimiterUser
	LoggerService   *utils.Logger
	MaxTokens       int
	RefillRate      time.Duration
	BlockDuration   time.Duration
	InactiveTimeout time.Duration
}

func RateLimiterMiddleware(rateLimiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RateLimiterHandler(next, rateLimiter)
	}
}

func RateLimiterHandler(next http.Handler, rateLimiter RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.RemoteAddr
		isExternal := strings.Contains(ipAddress, "[::1]")
		if !isExternal {
			ipAddress = strings.Split(ipAddress, ":")[0]
		}
		if _, exists := rateLimiter.Users[ipAddress]; exists {
			res := rateLimiter.Allow(ipAddress)
			if res {
				next.ServeHTTP(w, r)
				return
			} else {
				err := models.NewError(429, "Rate limiter", "Too many requests")
				utils.SetError(w, r, err)
				return
			}
		}
		newUser := NewRateLimiterUser(rateLimiter.MaxTokens, time.Now(), time.Unix(0, 0))
		rateLimiter.Users[ipAddress] = newUser
		rateLimiter.Allow(ipAddress)
		next.ServeHTTP(w, r)
	})
}

func NewRateLimiter(maxTokens int, refillRate time.Duration, blockDuration time.Duration, inactiveTimeout time.Duration, loggerService *utils.Logger) *RateLimiter {
	return &RateLimiter{
		Users:           map[string]*RateLimiterUser{},
		LoggerService:   loggerService,
		MaxTokens:       maxTokens,
		RefillRate:      refillRate,
		BlockDuration:   blockDuration,
		InactiveTimeout: inactiveTimeout,
	}
}

func (ra *RateLimiter) CleanWorker(ctx context.Context) {
	period := 60 * time.Second
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ra.LoggerService.Info("Started cleaning users in rate limiter")
			ra.Cleanup()
		case <-ctx.Done():
			ra.LoggerService.Info("Ended cleaning rate limiter")
			return
		}
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
	ra.LoggerService.Info("Removed users from rate limiter", removed)
}

func (ra *RateLimiter) Allow(ipAddress string) bool {
	now := time.Now()
	user := ra.Users[ipAddress]
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
	} else {
		user.BlockUntil = now.Add(ra.BlockDuration)
	}

	return false
}
