package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	redis "gopkg.in/redis.v5"

	"github.com/gorilla/handlers"
)

func main() {
	mux := http.NewServeMux()

	duration, err := time.ParseDuration("30s")
	if err != nil {
		log.Println(err)
	}

	mux.HandleFunc("/rate-limited", ratelimit(index, 300, duration))

	log.Println("Serving HTTP on port", "5000")
	http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, mux))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"response": "This is a rate limited response"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// middleware
func ratelimit(fn http.HandlerFunc, limit int64, per time.Duration) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		key := fmt.Sprintf("rate-limit/%s/", r.RemoteAddr)

		incr := client.Incr(key)
		// check if incr first value for key
		// add expire time
		if incr.Val() == 1 {
			client.Expire(key, per)
		}

		reset := strconv.Itoa(int(client.TTL(key).Val().Seconds()) + int(time.Now().UTC().Unix()))

		remaining := limit - incr.Val()
		if remaining < 0 {
			remaining = 0
		}

		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(limit)))
		w.Header().Set("X-RateLimit-Reset", reset)
		if incr.Val() > limit {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(429)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"data":  "You hit the rate limit",
				"error": "429"}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		fn(w, r)
	})
}
