package ratelimiter

import (
	"math"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	mu         sync.Mutex
	Capacity   float64
	Count      float64
	RefillRate float64
	LastRefill time.Time
}

func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
	return &TokenBucket{
		Capacity:   capacity,
		Count:      capacity,
		RefillRate: refillRate,
		LastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.LastRefill).Seconds()

	tb.LastRefill = now

	// Add tokens earned during elapsed time, capped at max capacity
	tb.Count = math.Min(tb.Capacity, tb.Count+(elapsed*tb.RefillRate))

	if tb.Count >= 1.0 {
		tb.Count--
		return true
	}

	return false
}

type UserLimiter struct {
	Capacity   float64
	RefillRate float64
	mu         sync.RWMutex
	buckets    map[string]*TokenBucket
}

func NewUserLimiter(capacity, refillrate float64) *UserLimiter {
	return &UserLimiter{
		Capacity:   capacity,
		RefillRate: refillrate,
		buckets:    make(map[string]*TokenBucket),
	}
}

func (ul *UserLimiter) GetBucket(userID string) *TokenBucket {
	ul.mu.RLock()
	bucket, exists := ul.buckets[userID]
	ul.mu.RUnlock()

	if exists {
		return bucket
	}

	ul.mu.Lock()
	defer ul.mu.Unlock()

	if bukcet, ok := ul.buckets[userID]; ok {
		return bukcet
	}

	ul.buckets[userID] = NewTokenBucket(ul.Capacity, ul.RefillRate)
	return ul.buckets[userID]
}

func (ul *UserLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = r.RemoteAddr
		}

		limiter := ul.GetBucket(userID)
		if !limiter.Allow() {
			w.Header().Set("Retry after", "1")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
