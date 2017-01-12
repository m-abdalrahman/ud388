package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	redis "gopkg.in/redis.v5"
)

func main() {
	mux := http.NewServeMux()

	duration, err := time.ParseDuration("60s")
	if err != nil {
		log.Println(err)
	}

	mux.HandleFunc("/catalog", ratelimit(getCatalog, 60, duration))

	log.Println("Serving HTTP on port", "5000")
	http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, mux))
}

func getCatalog(w http.ResponseWriter, r *http.Request) {
	items := []Item{}

	DB.Find(&items)

	if len(items) == 0 {
		item1 := Item{Name: "Pineapple", Price: "$2.50", Picture: "https://upload.wikimedia.org/wikipedia/commons/c/cb/Pineapple_and_cross_section.jpg", Description: "Organically Grown in Hawai'i"}
		DB.Create(&item1)
		item2 := Item{Name: "Carrots", Price: "$1.99", Picture: "http://media.mercola.com/assets/images/food-facts/carrot-fb.jpg", Description: "High in Vitamin A"}
		DB.Create(&item2)
		item3 := Item{Name: "Aluminum Foil", Price: "$3.50", Picture: "http://images.wisegeek.com/aluminum-foil.jpg", Description: "300 feet long"}
		DB.Create(&item3)
		item4 := Item{Name: "Eggs", Price: "$2.00", Picture: "http://whatsyourdeal.com/grocery-coupons/wp-content/uploads/2015/01/eggs.png", Description: "Farm Fresh Organic Eggs"}
		DB.Create(&item4)
		item5 := Item{Name: "Bananas", Price: "$2.15", Picture: "http://dreamatico.com/data_images/banana/banana-3.jpg", Description: "Fresh, delicious, and full of potassium"}
		DB.Create(&item5)

		DB.Find(&items)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ratelimit(fn http.HandlerFunc, limit int64, per time.Duration) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		key := fmt.Sprintf("rate-limit/%s/", r.RemoteAddr)

		incr := client.Incr(key)
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
