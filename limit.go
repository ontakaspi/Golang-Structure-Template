package main

import (
	"encoding/json"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Change the the map to hold values of the type visitor.
var visitors = make(map[string]*visitor)
var mu sync.Mutex

func getVisitor(indexKey string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[indexKey]
	if !exists {
		limiter := rate.NewLimiter(1, 10)
		// Include the current time when creating a new visitor.
		visitors[indexKey] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	mu.Lock()
	for indexKey, v := range visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(visitors, indexKey)
		}
	}
	mu.Unlock()
}

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var limiter *rate.Limiter
		jwtKey := r.Header.Get("Authorization")
		if jwtKey == "" {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			limiter = getVisitor(ip)
		} else {
			limiter = getVisitor(jwtKey)
		}

		if limiter.Allow() == false {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			type errorResp struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}

			json.NewEncoder(w).Encode(errorResp{
				Status:  "Too Many Requests",
				Message: "You have made too many requests in a given amount of time. Please try again later.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
