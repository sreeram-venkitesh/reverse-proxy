package proxy

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"time"
)

const (
	defaultExpiration = 5 * time.Minute
	purgeTime         = 10 * time.Minute
)

type CachedResponse struct {
	Body       []byte
	Headers    http.Header
	StatusCode int
}

func generateCacheKey(method, url string, headers http.Header) string {
	h := md5.New()
	io.WriteString(h, method)
	io.WriteString(h, url)

	if accept := headers.Get("Accept"); accept != "" {
		io.WriteString(h, accept)
	}

	return hex.EncodeToString(h.Sum(nil))
}

// Caching is done only for successful GET requests.
func shouldCache(method string, statusCode int) bool {
	if method != http.MethodGet {
		return false
	}

	return statusCode >= 200 && statusCode < 300
}
