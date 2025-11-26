package middleware

import (
	"context"
	"net/http"
	"strings"

	// "strings"
	"time"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
)

type RateLimiterUser struct {
	tokens     int
	lastRefill time.Time
	blockUntil time.Time
}

func NewRateLimiterUser(token int, lastRefill time.Time, blockUntil time.Time) *RateLimiterUser {
	return &RateLimiterUser{
		tokens:     token,
		lastRefill: lastRefill,
		blockUntil: blockUntil,
	}
}

type RateLimiter struct {
	users           map[string]*RateLimiterUser
	loggerService   utils.LoggerService
	maxTokens       int
	refillRate      time.Duration
	blockDuration   time.Duration
	inactiveTimeout time.Duration
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
		if _, exists := rateLimiter.users[ipAddress]; exists {
			res := rateLimiter.allow(ipAddress)
			if res {
				next.ServeHTTP(w, r)
				return
			} else {
				err := models.NewError(429, "Rate limiter", "Too many requests")
				response.SetError(w, r, err)
				return
			}
		}
		newUser := NewRateLimiterUser(rateLimiter.maxTokens, time.Now(), time.Unix(0, 0))
		rateLimiter.users[ipAddress] = newUser
		rateLimiter.allow(ipAddress)
		next.ServeHTTP(w, r)
	})
}

func NewRateLimiter(maxTokens int, refillRate time.Duration, blockDuration time.Duration, inactiveTimeout time.Duration, loggerService utils.LoggerService) *RateLimiter {
	return &RateLimiter{
		users:           map[string]*RateLimiterUser{},
		loggerService:   loggerService,
		maxTokens:       maxTokens,
		refillRate:      refillRate,
		blockDuration:   blockDuration,
		inactiveTimeout: inactiveTimeout,
	}
}

func (ra *RateLimiter) CleanWorker(ctx context.Context) {
	period := 60 * time.Second
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ra.loggerService.Info("Started cleaning users in rate limiter")
			ra.cleanup()
		case <-ctx.Done():
			ra.loggerService.Info("Ended cleaning rate limiter")
			return
		}
	}
}

func (ra *RateLimiter) cleanup() {
	now := time.Now()
	removeTime := now.Add(ra.inactiveTimeout)
	removed := 0
	for ipAddress, user := range ra.users {
		if removeTime.Before(user.lastRefill) {
			delete(ra.users, ipAddress)
			removed++
		}
	}
	ra.loggerService.Info("Removed users from rate limiter", removed)
}

func (ra *RateLimiter) allow(ipAddress string) bool {
	now := time.Now()
	user := ra.users[ipAddress]
	if now.Before(user.blockUntil) {
		return false
	}
	elapsed := now.Sub(user.lastRefill)
	newTokens := elapsed / ra.refillRate
	if newTokens > 0 {
		user.tokens += int(newTokens)
		if user.tokens > ra.maxTokens {
			user.tokens = ra.maxTokens
		}
		user.lastRefill = now
	}

	if user.tokens > 0 {
		user.tokens--
		return true
	} else {
		user.blockUntil = now.Add(ra.blockDuration)
	}

	return false
}
