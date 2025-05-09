package middleware

import (
	"net"
	"net/http"
	"strings"
)

var trustedProxies = []string{
	"10.0.0.0/8",     // Private network
	"172.16.0.0/12",  // Private network
	"192.168.0.0/16", // Private network
	"127.0.0.1/32",   // Localhost
	"100.64.0.0/10",  // Carrier-grade NAT
	"169.254.0.0/16", // Link-local
}

func isTrustedProxy(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, cidr := range trustedProxies {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func getClientIP(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}

	if !isTrustedProxy(remoteIP) {
		return remoteIP
	}

	headers := []string{"X-Forwarded-For", "X-Real-Ip", "CF-Connecting-IP"}

	for _, header := range headers {
		if ip := r.Header.Get(header); ip != "" {
			if header == "X-Forwarded-For" {
				ips := strings.Split(ip, ", ")
				for i := len(ips) - 1; i >= 0; i-- {
					ipStr := strings.TrimSpace(ips[i])
					if net.ParseIP(ipStr) != nil {
						if !isTrustedProxy(ipStr) {
							return ipStr
						}
					}
				}
			} else {
				if net.ParseIP(ip) != nil {
					return ip
				}
			}
		}
	}

	return remoteIP
}
