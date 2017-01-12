package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	googleAPIkey = ""
)

type JSONResults struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `lat:"location"`
				Lng float64 `lng:"location"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func main() {
	lat, lng := getGeocodeLocation("tokyo,japan")
	fmt.Println(lat, lng)
}

// Use Google Maps to convert a location into Latitute/Longitute coordinates
// FORMAT: https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&key=API_KEY
func getGeocodeLocation(inputString string) (lat, lng float64) {
	locationString := strings.Replace(inputString, " ", "+", -1)

	// urlAddress := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
	// locationString, googleAPIkey)
	urlAddress := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s",
		locationString)

	resp, err := http.Get(urlAddress)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	var result JSONResults

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	lat = result.Results[0].Geometry.Location.Lat
	lng = result.Results[0].Geometry.Location.Lng

	return
}
