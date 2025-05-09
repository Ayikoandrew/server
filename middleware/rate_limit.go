package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type RateLimitConfig struct {
	RequestPerMinute       int
	GlobalRequestPerMinute int
	CleanUpInterval        time.Duration
	IPExpiryInterval       time.Duration
}

var DefaultConfig = RateLimitConfig{
	RequestPerMinute:       20,
	GlobalRequestPerMinute: 200,
	CleanUpInterval:        5 * time.Minute,
	IPExpiryInterval:       30 * time.Minute,
}

type TokenBucket struct {
	capacity     int
	tokens       int
	refillRate   int
	lastRefill   time.Time
	lastAccessed time.Time
	mu           sync.Mutex
	ipBuckets    map[string]*TokenBucket
	globalBucket *TokenBucket
	config       RateLimitConfig
}

func NewTokenBucket(requestPerMinute, refillRate int) *TokenBucket {
	return NewTokenBucketWithConfig(RateLimitConfig{
		RequestPerMinute:       requestPerMinute,
		GlobalRequestPerMinute: requestPerMinute * 10,
		CleanUpInterval:        5 * time.Minute,
		IPExpiryInterval:       30 * time.Minute,
	})
}

func NewTokenBucketWithConfig(config RateLimitConfig) *TokenBucket {
	now := time.Now()
	tb := &TokenBucket{
		capacity:     config.RequestPerMinute,
		tokens:       config.RequestPerMinute,
		refillRate:   config.RequestPerMinute / 60,
		lastRefill:   now,
		lastAccessed: now,
		ipBuckets:    make(map[string]*TokenBucket),
		config:       config,
	}

	tb.globalBucket = &TokenBucket{
		capacity:     config.RequestPerMinute,
		tokens:       config.RequestPerMinute,
		refillRate:   config.RequestPerMinute / 60,
		lastRefill:   now,
		lastAccessed: now,
		config:       config,
	}

	go tb.startCleaner()
	return tb
}

func (t *TokenBucket) startCleaner() {
	ticker := time.NewTicker(t.globalBucket.config.CleanUpInterval)
	defer ticker.Stop()

	for range ticker.C {
		t.cleanup()
	}
}

func (t *TokenBucket) cleanup() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	threshold := now.Add(-t.config.IPExpiryInterval)

	for ip, bucket := range t.ipBuckets {
		if bucket.lastAccessed.Before(threshold) {
			delete(t.ipBuckets, ip)
			log.Printf("Cleaned up rate limit bucket for IP: %s", ip)
		}
	}
}

func (t *TokenBucket) Allow(ip string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()

	t.globalBucket.refill()
	if t.globalBucket.tokens <= 0 {
		log.Printf("Rate limit exceeded: global limit reached")
		return false
	}

	if _, exists := t.ipBuckets[ip]; !exists {
		t.ipBuckets[ip] = &TokenBucket{
			capacity:     t.capacity,
			tokens:       t.capacity,
			refillRate:   t.refillRate,
			lastRefill:   now,
			lastAccessed: now,
		}
		log.Printf("New client connected %s\n", ip)
	}

	ipBucket := t.ipBuckets[ip]
	ipBucket.lastAccessed = now
	ipBucket.refill()

	if ipBucket.tokens > 0 && t.globalBucket.tokens > 0 {
		ipBucket.tokens--
		t.globalBucket.tokens--
		ipBucket.lastAccessed = now
		return true
	}

	if ipBucket.tokens <= 0 {
		log.Printf("Rate limit exceeded for IP: %s", ip)
	}

	return false
}

func (t *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(t.lastAccessed)
	tokensToAdd := int(elapsed.Seconds()) * t.refillRate

	if tokensToAdd > 0 {
		t.tokens += tokensToAdd
		if t.tokens > t.capacity {
			t.tokens = t.capacity
		}
		t.lastRefill = now
	}
}

var (
	globalLimiter *TokenBucket
	once          sync.Once
)

func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	once.Do(func() {
		globalLimiter = NewTokenBucketWithConfig(DefaultConfig)
		log.Printf("Rate limit initialized %d requests per minute per IP, %d global",
			DefaultConfig.RequestPerMinute,
			DefaultConfig.GlobalRequestPerMinute,
		)
	})

	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		if !globalLimiter.Allow(clientIP) {
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func RateLimitMiddlewareTokenBucket(next http.HandlerFunc) http.HandlerFunc {
	return rateLimitMiddleware(next)
}
