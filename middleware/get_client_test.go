package middleware

import (
	"testing"
)

func TestIsTrustedProxy(t *testing.T) {
	trustedIPs := []string{
		"10.1.2.3",    // Within 10.0.0.0/8
		"172.16.0.1",  // Within 172.16.0.0/12
		"192.168.1.1", // Within 192.168.0.0/16
		"127.0.0.1",   // Localhost
		"100.64.0.1",  // Within 100.64.0.0/10
		"169.254.1.1", // Link-local
	}

	for _, ip := range trustedIPs {
		if !isTrustedProxy(ip) {
			t.Errorf("Expected %s to be trusted, but it was not", ip)
		}
	}

	untrustedIPs := []string{
		"",
		"8.8.8.8",         // Google DNS
		"1.1.1.1",         // Cloudflare DNS
		"203.0.113.1",     // TEST-NET-3
		"198.51.100.1",    // TEST-NET-2
		"192.0.2.1",       // TEST-NET-1
		"invalid-ip",      // Invalid format
		"256.256.256.256", // Invalid IP
		"2001:db8::1",     // IPv6 documentation
	}

	for _, ip := range untrustedIPs {
		if isTrustedProxy(ip) {
			t.Errorf("Expected %s to be untrusted, but it was trusted", ip)
		}
	}
}
