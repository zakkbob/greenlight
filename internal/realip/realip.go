// adapted from http://github.com/tomasen/realip
// original license: MIT

// WARNING: this is a bad system, X-Forwarded-For and X-Real-IP should only be used from
// trusted proxies, clients can easily spoof these headers

package realip

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type ErrInvalidAddress struct {
	Addr string
}

func (e ErrInvalidAddress) Error() string {
	return fmt.Sprintf("invalid address: %s", e.Addr)
}

func isAllowedIP(ip net.IP) bool {
	return ip.IsGlobalUnicast() &&
		!ip.IsPrivate() &&
		!ip.IsLinkLocalUnicast() &&
		!ip.IsLoopback() &&
		!ip.IsMulticast() &&
		!ip.IsUnspecified()
}

func stripPort(addr string) (string, error) {
	if !strings.ContainsRune(addr, ':') {
		return addr, nil
	}

	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", fmt.Errorf("stripping port: %w", err)
	}
	return ip, nil
}

// FromRequest return client's real public IP address from http request headers.
func FromRequest(r *http.Request) (string, error) {
	xRealIP := r.Header.Get("X-Real-Ip")
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	for addr := range strings.SplitSeq(xForwardedFor, ",") {
		addr = strings.TrimSpace(addr)
		ip := net.ParseIP(addr)
		if ip != nil && isAllowedIP(ip) {
			return addr, nil
		}
	}

	if ip := net.ParseIP(xRealIP); ip != nil {
		return xRealIP, nil
	}

	addr, err := stripPort(r.RemoteAddr)
	if err != nil {
		return "", ErrInvalidAddress{Addr: addr}
	}
	if ip := net.ParseIP(addr); ip == nil {
		return "", ErrInvalidAddress{Addr: addr}
	}
	return addr, nil
}
