package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type item struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Price       string `json:"price"`
	Description string `json:"description"`
}

func main() {
	url := "http://localhost:5000/catalog"

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
		defer resp.Body.Close()

		var result interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Println(err)
		}

		switch result.(type) {
		case []interface{}:
			fmt.Println("Number of Requests:", requests+1)
			value := result.([]interface{})
			// https://github.com/golang/go/wiki/InterfaceSlice
			// mapValue := value[0].(map[string]interface{})
			// fmt.Println(mapValue["id"])
			fmt.Println(value)
		case map[string]interface{}:
			errMsg := result.(map[string]interface{})
			fmt.Printf("Error #%s : %s\n", errMsg["error"], errMsg["data"])
			fmt.Println("Hit rate limit. Waiting 5 seconds and trying again...")
			time.Sleep(5 * time.Second)
			SendRequests(url, reqPerMinute)
		}

		requests = requests + 1
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
