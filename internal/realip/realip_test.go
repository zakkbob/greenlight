// adapted from http://github.com/tomasen/realip
// original license: MIT

package realip_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/zakkbob/greenlight/internal/realip"
)

// dodgy tests
func TestRealIP(t *testing.T) {
	newRequest := func(remoteAddr, xRealIP string, xForwardedFor ...string) *http.Request {
		h := http.Header{}
		h.Set("X-Real-IP", xRealIP)
		h.Set("X-Forwarded-For", strings.Join(xForwardedFor, ","))

		return &http.Request{
			RemoteAddr: remoteAddr,
			Header:     h,
		}
	}

	publicAddr1 := "144.12.54.87"
	publicAddr2 := "119.14.55.11"
	localAddr := "127.0.0.0"

	testData := []struct {
		name     string
		request  *http.Request
		expected string
	}{
		{
			name:     "No header",
			request:  newRequest(publicAddr1, ""),
			expected: publicAddr1,
		}, {
			name:     "Has X-Forwarded-For",
			request:  newRequest("", "", publicAddr1),
			expected: publicAddr1,
		}, {
			name:     "Has multiple X-Forwarded-For",
			request:  newRequest("", "", localAddr, publicAddr1, publicAddr2),
			expected: publicAddr1,
		}, {
			name:     "Has X-Real-IP",
			request:  newRequest("", publicAddr1),
			expected: publicAddr1,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := realip.FromRequest(tt.request)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expected != actual {
				t.Errorf("%s: expected %s but get %s", tt.name, tt.expected, actual)
			}
		})
	}
}
