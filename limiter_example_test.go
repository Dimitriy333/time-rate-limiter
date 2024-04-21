package time_rate_limiter

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func LimiterExampleTest() {
	limiterCleanUpInterval := 5 * time.Minute
	messageLimiter := NewRateLimiter(limiterCleanUpInterval)
	transactionLimiter := NewRateLimiter(limiterCleanUpInterval)
	ipLimiter := NewRateLimiter(limiterCleanUpInterval)

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User-Id")
		if user == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		// A user can send no more than 5 messages per second
		if messageLimiter.Limit(user, 5, time.Second) { // A user can send no more than 5 messages per second
			_, err := fmt.Fprintf(w, "Message sent successfully")
			if err != nil {
				log.Println(err)
			}
		} else {
			http.Error(w, "Rate limit exceeded for messages", http.StatusTooManyRequests)
		}
	})

	http.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User-Id")
		if user == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		// A user can have no more than 3 failed card transactions per day
		if transactionLimiter.Limit(user, 3, 24*time.Hour) {
			_, err := fmt.Fprintf(w, "Transaction processed successfully")
			if err != nil {
				log.Println(err)
			}
		} else {
			http.Error(w, "Rate limit exceeded for card failures", http.StatusTooManyRequests)
		}
	})

	http.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-User-Ip")
		if ip == "" {
			http.Error(w, "User IP is required", http.StatusBadRequest)
			return
		}

		// One IP address can send no more than 10,000 requests per minute
		if ipLimiter.Limit(ip, 10000, time.Minute) {
			_, err := fmt.Fprintf(w, "Request processed successfully")
			if err != nil {
				log.Println(err)
			}
		} else {
			http.Error(w, "Rate limit exceeded for IP", http.StatusTooManyRequests)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}

	// dispose resources at the end
	messageLimiter.Dispose()
	transactionLimiter.Dispose()
	ipLimiter.Dispose()
}
