package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	url := "http://localhost:5000/rate-limited"

	var reqPerMinute int

	fmt.Print("Please specify the number of requests per minute: ")
	fmt.Scan(&reqPerMinute)

	fmt.Println("Sending Requests...")
	SendRequests(url, reqPerMinute)
}

func SendRequests(url string, reqPerMinute int) {
	interval := 60 / reqPerMinute
	requests := 0
	for requests < reqPerMinute {
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		result := map[string]string{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Println(err)
		}

		if _, ok := result["error"]; ok {
			fmt.Printf("Error #%s : %s\n", result["error"], result["data"])
			fmt.Println("Hit rate limit. Waiting 5 seconds and trying again...")
			time.Sleep(5 * time.Second)
			SendRequests(url, reqPerMinute)
		} else {
			fmt.Println("Number of Requests:", requests+1)
			fmt.Println(result["response"])
		}

		requests = requests + 1
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
