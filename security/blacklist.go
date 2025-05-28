package security

import (
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Blacklist struct {
	tokenBlacklist map[string]time.Time
	mu             sync.RWMutex
}

func NewBlacklist() *Blacklist {
	return &Blacklist{
		tokenBlacklist: make(map[string]time.Time),
	}
}

func (b *Blacklist) AddTokenToBlacklist(refreshToken string, expiry time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokenBlacklist[refreshToken] = expiry
}

func (b *Blacklist) IsTokenBlacklisted(token string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	expiry, exists := b.tokenBlacklist[token]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		b.mu.Lock()
		delete(b.tokenBlacklist, token)
		b.mu.Unlock()
		return false
	}
	return true
}

func (b *Blacklist) CleanUpBlacklist() {
	ticker := time.NewTicker(7 * 24 * 60 * 60)

	for range ticker.C {
		for token, expiry := range b.tokenBlacklist {
			b.mu.Lock()
			now := time.Now()

			if now.After(expiry) {
				delete(b.tokenBlacklist, token)
				slog.Info("Expired token deleted")
			}
			b.mu.Unlock()
		}
	}
}

var (
	blacklist *Blacklist
	once      sync.Once
)

func blacklistMiddleware(next http.Handler) http.HandlerFunc {
	once.Do(func() {
		blacklist := NewBlacklist()
		log.Printf("Blacklist %s", blacklist.tokenBlacklist)
	})
	return func(w http.ResponseWriter, r *http.Request) {
		refresh, err := r.Cookie("refresh_token")
		if err != nil {
			slog.Error("cookie does not exist", "err", err)
			return
		}

		blacklist.AddTokenToBlacklist(refresh.Value, refresh.Expires)

		if time.Now().After(refresh.Expires) {
			if blacklist.IsTokenBlacklisted(refresh.Value) {
				blacklist.CleanUpBlacklist()
			}
		}

		next.ServeHTTP(w, r)
	}
}

func BlacklistMiddleware(next http.Handler) http.HandlerFunc {
	return blacklistMiddleware(next)
}
