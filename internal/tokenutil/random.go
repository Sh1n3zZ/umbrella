// Package tokenutil provides cryptographically random string helpers used
// across OAuth flows (authorization codes, verification tokens, etc.).
package tokenutil

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

const (
	// randomEntropyMinBytes is the smallest random byte count accepted before encoding.
	randomEntropyMinBytes = 1
	// randomEntropyMaxBytes is the largest random byte count accepted before encoding.
	randomEntropyMaxBytes = 4096
)

// GenerateURLSafeRandomString reads n cryptographically random bytes and returns
// base64.URLEncoding of that buffer (RFC 4648 URL-safe alphabet: '-' and '_' instead
// of '+' and '/'). Output uses only [A-Za-z0-9-_=], which is a subset of RFC 6749
// VSCHAR for values like "state" (Appendix A.5), "code" (A.11), and tokens (A.12, A.17).
// Returned string length is base64.URLEncoding.EncodedLen(n).
func GenerateURLSafeRandomString(n int) (string, error) {
	if n < randomEntropyMinBytes || n > randomEntropyMaxBytes {
		return "", fmt.Errorf("tokenutil: random entropy byte count must be between %d and %d", randomEntropyMinBytes, randomEntropyMaxBytes)
	}

	raw := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, raw); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(raw), nil
}
