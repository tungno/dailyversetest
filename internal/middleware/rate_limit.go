/**
 *  RateLimitMiddleware provides middleware to limit the number of requests per client IP.
 *  This implementation uses a token bucket algorithm provided by the `golang.org/x/time/rate`
 *  package to enforce rate limits and maintain fairness among clients.
 *
 *  @file       rate_limit.go
 *  @package    middleware
 *
 *  @properties
 *  - clients (map[string]*client) - A map storing rate limiters for each client IP.
 *  - mutex (sync.Mutex)           - A mutex to ensure thread-safe access to the clients map.
 *  - rateLimit (rate.Limit)       - The rate of requests allowed per time period.
 *  - burst (int)                  - The maximum burst size of requests allowed.
 *  - cleanupInterval (time.Duration) - The interval to clean up inactive clients.
 *
 *  @struct   client
 *  - limiter (*rate.Limiter) - A token bucket rate limiter for the client.
 *  - lastSeen (time.Time)    - The last time this client was active.
 *
 *  @methods
 *  - RateLimitMiddleware(next)       - Middleware to enforce rate limiting on requests.
 *  - getIP(r)                        - Extracts the client's IP address from the HTTP request.
 *  - cleanupClients()                - Periodically removes inactive clients from the map.
 *
 *  @behavior
 *  - Enforces a maximum of 5 requests per hour per client IP.
 *  - Allows bursts of up to 5 requests within the defined time period.
 *  - Returns a 429 Too Many Requests error if the client exceeds the rate limit.
 *  - Automatically cleans up clients that have been inactive for a specified duration.
 *
 *  @example
 *  ```
 *  func main() {
 *      mux := http.NewServeMux()
 *      mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
 *          w.Write([]byte("Hello, world!"))
 *      })
 *
 *      handler := middleware.RateLimitMiddleware(mux)
 *      http.ListenAndServe(":8080", handler)
 *  }
 *  ```
 *
 *  @dependencies
 *  - "golang.org/x/time/rate": Provides token bucket rate limiting.
 *  - sync.Mutex: Ensures thread-safe access to shared resources.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

// client represents a single client's rate limiter and last activity.
type client struct {
	limiter  *rate.Limiter // Rate limiter for the client.
	lastSeen time.Time     // Timestamp of the client's last request.
}

var (
	clients         = make(map[string]*client)  // Map of client IPs to rate limiters.
	mutex           sync.Mutex                  // Mutex for thread-safe map access.
	rateLimit       = rate.Every(time.Hour / 5) // 5 requests per hour.
	burst           = 5                         // Burst size: maximum number of requests in quick succession.
	cleanupInterval = time.Minute * 10          // Interval to clean up inactive clients.
)

// RateLimitMiddleware limits the number of requests per client.
func RateLimitMiddleware(next http.Handler) http.Handler {
	// Start the client cleanup goroutine.
	go cleanupClients()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the client's IP address.
		ip := getIP(r)

		mutex.Lock()
		// Retrieve or initialize the client's rate limiter.
		c, exists := clients[ip]
		if !exists {
			limiter := rate.NewLimiter(rateLimit, burst)
			clients[ip] = &client{limiter: limiter, lastSeen: time.Now()}
			c = clients[ip]
		}
		// Update the client's last seen timestamp.
		c.lastSeen = time.Now()
		mutex.Unlock()

		// Enforce the rate limit.
		if !c.limiter.Allow() {
			http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}

		// Proceed to the next handler.
		next.ServeHTTP(w, r)
	})
}

// getIP extracts the client's real IP address from the request headers or RemoteAddr.
func getIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs; use the first IP.
		return xff
	}
	return r.RemoteAddr
}

// cleanupClients periodically removes inactive clients from the map.
func cleanupClients() {
	for {
		time.Sleep(cleanupInterval)
		mutex.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > cleanupInterval {
				delete(clients, ip)
			}
		}
		mutex.Unlock()
	}
}
