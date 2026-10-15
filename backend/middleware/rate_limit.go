package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mini-project/backend/database"

	"github.com/redis/go-redis/v9"
)

var rateLimitScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
    redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
local ttl = redis.call("PTTL", KEYS[1])
return {current, ttl}
`)

func RateLimit(limit int, window time.Duration, keyPrefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if database.Redis == nil {
				writeRateLimitError(w, http.StatusServiceUnavailable, "service temporarily unavailable")
				return
			}

			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				host = r.RemoteAddr
			}

			key := keyPrefix + ":" + strings.TrimSpace(host)

			ctx := r.Context()
			result, err := rateLimitScript.Run(
				ctx,
				database.Redis,
				[]string{key},
				window.Milliseconds(),
			).Slice()

			if err != nil || len(result) < 2 {
				writeRateLimitError(w, http.StatusServiceUnavailable, "service temporarily unavailable")
				return
			}

			var count int64
			switch value := result[0].(type) {
			case int64:
				count = value
			case string:
				count, _ = strconv.ParseInt(value, 10, 64)
			case []byte:
				count, _ = strconv.ParseInt(string(value), 10, 64)
			}

			var ttl int64
			switch value := result[1].(type) {
			case int64:
				ttl = value
			case string:
				ttl, _ = strconv.ParseInt(value, 10, 64)
			case []byte:
				ttl, _ = strconv.ParseInt(string(value), 10, 64)
			}

			if count > int64(limit) {
				retryAfter := int64(time.Duration(ttl) * time.Millisecond / time.Second)
				if retryAfter < 1 {
					retryAfter = 1
				}
				w.Header().Set("Retry-After", strconv.FormatInt(retryAfter, 10))
				w.Header().Set("RateLimit-Limit", strconv.Itoa(limit))
				w.Header().Set("RateLimit-Remaining", "0")
				writeRateLimitError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}

			remaining := int64(limit) - count
			w.Header().Set("RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("RateLimit-Remaining", strconv.FormatInt(remaining, 10))
			w.Header().Set("RateLimit-Reset", strconv.FormatInt(ttl/1000, 10))

			next.ServeHTTP(w, r)
		})
	}
}

func writeRateLimitError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
