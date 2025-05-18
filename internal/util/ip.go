package util

import (
	"net"
	"net/http"
	"strings"
)

func ReadUserIP(r *http.Request) (net.IP, error) {
	// Check X-Real-IP header
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return net.ParseIP(ip), nil
	}

	// Check X-Forwarded-For header (may contain multiple IPs)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		return net.ParseIP(ip), nil
	}

	// Fallback to RemoteAddr (host:port)
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(host), nil
}
