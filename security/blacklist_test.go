package security

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"
)

// github.com/stretchr/testify
func TestAddTokenToBlacklist(t *testing.T) {
	table := []string{}

	b := NewBlacklist()
	t.Run("Adding token to blacklist", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			token := randomBytes(16)
			t := time.Time.Add(time.Time{}, 2*time.Second)
			b.AddTokenToBlacklist(token, t)
			table = append(table, token)
		}
	})

	t.Run("Checking token", func(t *testing.T) {
		for _, tk := range table {

			if b.IsTokenBlacklisted(tk) {
				t.Errorf("Expected token %s to be blacklisted", tk)
			}

		}
	})

	// Test takes more than 30s
	t.Run("Clean up blacklist", func(t *testing.T) {
		b.CleanUpBlacklist()
	})
}

func randomBytes(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return hex.EncodeToString(b)
}
