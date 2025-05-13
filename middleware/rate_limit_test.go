package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
)

// NOTE:
// Some of the tests are repeated and they are my playground for learning to write better tests.
// I got to learn :)

func TestTokenBucket_Allow(t *testing.T) {
	config := RateLimitConfig{
		RequestPerMinute:       5,
		GlobalRequestPerMinute: 15,
		CleanUpInterval:        1 * time.Minute,
		IPExpiryInterval:       5 * time.Minute,
	}

	t.Run("single IP within limit", func(t *testing.T) {
		tb := NewTokenBucketWithConfig(config, context.Background())

		for i := 0; i < 5; i++ {
			if !tb.Allow("192.168.1.1") {
				t.Errorf("Expected request %d to be allowed", i+1)
			}
		}
	})

	t.Run("single IP over limit", func(t *testing.T) {
		tb := NewTokenBucketWithConfig(config, context.Background())

		for i := 0; i < 5; i++ {
			tb.Allow("192.168.1.1")
		}

		if tb.Allow("192.168.1.1") {
			t.Error("Expected request to be denied (over limit)")
		}
	})

	t.Run("multiple IPs within global limit", func(t *testing.T) {
		tb := NewTokenBucketWithConfig(config, context.Background())

		ips := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}
		for _, ip := range ips {
			for i := 0; i < 5; i++ {
				if !tb.Allow(ip) {
					t.Errorf("Expected request from %s to be allowed", ip)
				}
			}
		}
	})

	t.Run("global limit exceeded", func(t *testing.T) {
		specialConfig := RateLimitConfig{
			RequestPerMinute:       5,
			GlobalRequestPerMinute: 1,
			CleanUpInterval:        1 * time.Minute,
			IPExpiryInterval:       5 * time.Minute,
		}
		tb := NewTokenBucketWithConfig(specialConfig, context.Background())

		tb.Allow("10.0.0.4")

		if tb.Allow("10.0.0.5") {
			t.Error("Expected request to be denied (global limit exceeded)")
		}
	})
}

func TestTokenBucket_Refill(t *testing.T) {
	config := RateLimitConfig{
		RequestPerMinute:       60,
		GlobalRequestPerMinute: 120,
		CleanUpInterval:        time.Minute,
		IPExpiryInterval:       5 * time.Minute,
	}
	tb := NewTokenBucketWithConfig(config, context.Background())

	for i := 0; i < 60; i++ {
		tb.Allow("192.168.1.1")
	}

	if tb.Allow("192.168.1.1") {
		t.Error("Expected to be rate limited after exhausting tokens")
	}

	time.Sleep(1100 * time.Millisecond)

	if !tb.Allow("192.168.1.1") {
		t.Error("Expected token to be available after refill period")
	}
}

func TestTokenBucket_Cleanup(t *testing.T) {
	config := RateLimitConfig{
		RequestPerMinute:       10,
		GlobalRequestPerMinute: 20,
		CleanUpInterval:        100 * time.Millisecond,
		IPExpiryInterval:       200 * time.Millisecond,
	}
	tb := NewTokenBucketWithConfig(config, context.Background())

	tb.Allow("10.0.0.1")
	tb.Allow("10.0.0.2")
	tb.Allow("10.0.0.3")

	if len(tb.ipBuckets) != 3 {
		t.Errorf("Expected 2 IP buckets, got %d", len(tb.ipBuckets))
	}

	time.Sleep(300 * time.Millisecond)

	tb.mu.Lock()
	defer tb.mu.Unlock()
	if len(tb.ipBuckets) != 0 {
		t.Errorf("Expected IP buckets to be cleaned up, got %d", len(tb.ipBuckets))
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	t.Run("single IP within limit", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		testConfig := RateLimitConfig{
			RequestPerMinute:       2,
			GlobalRequestPerMinute: 5,
			CleanUpInterval:        time.Minute,
			IPExpiryInterval:       5 * time.Minute,
		}

		testLimiter := NewTokenBucketWithConfig(testConfig, context.Background())

		testMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				clientIP := getClientIP(r)

				if !testLimiter.Allow(clientIP) {
					http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(w, r)
			})
		}

		ts := httptest.NewServer(testMiddleware(handler))
		defer ts.Close()

		client := ts.Client()

		for i := 0; i < 2; i++ {
			resp, err := client.Get(ts.URL)
			if err != nil {
				t.Fatal(err)
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK, got %d", resp.StatusCode)
			}
		}
	})

	t.Run("single IP over limit", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		testConfig := RateLimitConfig{
			RequestPerMinute:       2,
			GlobalRequestPerMinute: 5,
			CleanUpInterval:        time.Minute,
			IPExpiryInterval:       5 * time.Minute,
		}

		testLimiter := NewTokenBucketWithConfig(testConfig, context.Background())

		testMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				clientIP := getClientIP(r)

				if !testLimiter.Allow(clientIP) {
					http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(w, r)
			})
		}

		ts := httptest.NewServer(testMiddleware(handler))
		defer ts.Close()

		client := ts.Client()

		for i := 0; i < 2; i++ {
			resp, err := client.Get(ts.URL)
			if err != nil {
				t.Fatal(err)
			}
			resp.Body.Close()
		}

		for i := 0; i < 5; i++ {
			resp, err := client.Get(ts.URL)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("Expected status TooManyRequests, got %d", resp.StatusCode)
			}
		}
	})

	t.Run("different IPs within global limit", func(t *testing.T) {
		// This test would need multiple clients with different IPs to properly test
		// Skipping the implementation here
		t.Skip("This test requires a mock for IP address extraction")
	})
}

func TestRateLimitConcurrency(t *testing.T) {
	config := RateLimitConfig{
		RequestPerMinute:       100,
		GlobalRequestPerMinute: 1000,
		CleanUpInterval:        time.Minute,
		IPExpiryInterval:       5 * time.Minute,
	}
	tb := NewTokenBucketWithConfig(config, context.Background())

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if tb.Allow("10.0.0.1") {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if successCount != 100 {
		t.Errorf("Expected exactly 100 successful requests, got %d", successCount)
	}
}

func TestRateLimitWithMultipleIps(t *testing.T) {
	testConfig := RateLimitConfig{
		RequestPerMinute:       2,
		GlobalRequestPerMinute: 5,
		CleanUpInterval:        time.Minute,
		IPExpiryInterval:       time.Minute,
	}

	ipCounter := 0
	mockIPs := []string{
		"192.168.1.1",
		"192.168.1.2",
		"192.168.1.3",
		"192.168.1.4",
		"192.168.1.5",
		"192.168.1.6",
		"192.168.1.1",
	}

	mockGetClientIP := func(_ *http.Request) string {
		ip := mockIPs[ipCounter%len(mockIPs)]
		ipCounter++
		return ip
	}

	testLimiter := NewTokenBucketWithConfig(testConfig, context.TODO())

	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := mockGetClientIP(r)

			if !testLimiter.Allow(ip) {
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(testMiddleware(handler))
	defer ts.Close()

	client := ts.Client()

	responses := make([]struct {
		IP         string
		StatusCode int
	}, 6)

	ipCounter = 0

	for i := 0; i < 6; i++ {
		ip := mockIPs[i%len(mockIPs)]

		resp, err := client.Get(ts.URL)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}

		responses[i] = struct {
			IP         string
			StatusCode int
		}{
			IP:         ip,
			StatusCode: resp.StatusCode,
		}
		resp.Body.Close()
	}

	for i := 0; i < 5; i++ {
		if responses[i].StatusCode != http.StatusOK {
			t.Errorf("Request %d (IP: %s) expected status code 200, got %d", i, responses[i].IP, responses[i].StatusCode)
		}
	}

	if responses[5].StatusCode != http.StatusTooManyRequests {
		t.Errorf("Request 5 (IP: %s) expected status 429, got %d",
			responses[5].IP, responses[5].StatusCode)
	}
}

func TestRateLimitPerIPLimits(t *testing.T) {

	testConfig := RateLimitConfig{
		RequestPerMinute:       2,
		GlobalRequestPerMinute: 10,
		CleanUpInterval:        time.Minute,
		IPExpiryInterval:       5 * time.Minute,
	}

	mockIPs := []string{"10.0.0.1", "10.0.0.2"}

	mockGetClientIP := func(r *http.Request) string {
		ipIndex, _ := strconv.Atoi(r.Header.Get("X-Test-IP-Index"))
		return mockIPs[ipIndex%len(mockIPs)]
	}

	testLimiter := NewTokenBucketWithConfig(testConfig, context.Background())

	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := mockGetClientIP(r)

			if !testLimiter.Allow(clientIP) {
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(testMiddleware(handler))
	defer ts.Close()

	client := ts.Client()

	makeRequest := func(ipIndex int) int {
		req, _ := http.NewRequest("GET", ts.URL, nil)
		req.Header.Set("X-Test-IP-Index", strconv.Itoa(ipIndex))

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		return resp.StatusCode
	}

	// Test IP-specific rate limits

	// First two requests for IP1 should succeed
	if status := makeRequest(0); status != http.StatusOK {
		t.Errorf("First request for IP1 expected status 200, got %d", status)
	}
	if status := makeRequest(0); status != http.StatusOK {
		t.Errorf("Second request for IP1 expected status 200, got %d", status)
	}

	// Third request for IP1 should be limited
	if status := makeRequest(0); status != http.StatusTooManyRequests {
		t.Errorf("Third request for IP1 expected status 429, got %d", status)
	}

	// But requests for IP2 should still succeed
	if status := makeRequest(1); status != http.StatusOK {
		t.Errorf("First request for IP2 expected status 200, got %d", status)
	}
	if status := makeRequest(1); status != http.StatusOK {
		t.Errorf("Second request for IP2 expected status 200, got %d", status)
	}

	// And IP2 should be limited after its limit
	if status := makeRequest(1); status != http.StatusTooManyRequests {
		t.Errorf("Third request for IP2 expected status 429, got %d", status)
	}
}
