package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	googleAPIkey = ""

	foursquareClientID     = ""
	foursquareClientSecret = ""
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

type JSONRestaurantResponse struct {
	Response struct {
		Venues []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Location struct {
				FormattedAddress []string `json:"formattedAddress"`
			} `json:"location"`
		} `json:"venues"`
	} `json:"response"`
}

type JSONImgResponse struct {
	Response struct {
		Photos struct {
			Items []struct {
				Prefix string `json:"prefix"`
				Suffix string `json:"suffix"`
			} `json:"items"`
		} `json:"photos"`
	} `json:"response"`
}

// Use Google Maps to convert a location into Latitute/Longitute coordinates
// FORMAT: https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&key=API_KEY
func GetGeocodeLocation(inputString string) (lat, lng float64) {
	locationString := strings.Replace(inputString, " ", "+", -1)
	// urlAddress := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
	// 	locationString, googleAPIKey)
	urlAddress := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s", locationString)

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

func FindARestaurant(mealType, location string) (*Restaurant, error) {
	//1. Use getGeocodeLocation to get the latitude and longitude coordinates of the location string.
	lat, lng := GetGeocodeLocation(location)
	//2.  Use foursquare API to find a nearby restaurant with the latitude, longitude, and mealType strings.
	//HINT: format for url will be something like https://api.foursquare.com/v2/venues/search?client_id=CLIENT_ID&client_secret=CLIENT_SECRET&v=20130815&ll=40.7,-74&query=sushi
	urlAddress := fmt.Sprintf("https://api.foursquare.com/v2/venues/search?client_id=%s&client_secret=%s&v=20130815&ll=%s,%s&query=%s",
		foursquareClientID, foursquareClientSecret,
		strconv.FormatFloat(lat, 'f', -1, 64), strconv.FormatFloat(lng, 'f', -1, 64), mealType)

	resp, err := http.Get(urlAddress)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	var result JSONRestaurantResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Decode:", err)
	}

	if len(result.Response.Venues) != 0 {
		//3.  Grab the first restaurant
		// restaurant name
		restaurantName := result.Response.Venues[0].Name
		//restaurant address
		restaurantAddress := strings.Join(result.Response.Venues[0].Location.FormattedAddress, " ")
		//4.  Get a  300x300 picture of the restaurant using the venue_id (you can change this by altering the 300x300 value in the URL or replacing it with 'orginal' to get the original picture
		imgURL := fmt.Sprintf("https://api.foursquare.com/v2/venues/%s/photos?client_id=%s&v=20150603&client_secret=%s",
			result.Response.Venues[0].ID, foursquareClientID, foursquareClientSecret)

		resp, err := http.Get(imgURL)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		var imgResponse JSONImgResponse
		if err := json.NewDecoder(resp.Body).Decode(&imgResponse); err != nil {
			log.Println("Decode:", err)
		}

		var imageURL string
		//5.  Grab the first image
		if len(imgResponse.Response.Photos.Items) != 0 {
			prefix := imgResponse.Response.Photos.Items[0].Prefix
			suffix := imgResponse.Response.Photos.Items[0].Suffix
			imageURL = prefix + "300x300" + suffix
		} else {
			//6.  if no image available, insert default image url
			imageURL = "https://s-media-cache-ak0.pinimg.com/originals/02/ca/48/02ca48c68c5e31572f9e690815517032.jpg"
		}

		return &Restaurant{
			Name:    restaurantName,
			Address: restaurantAddress,
			Image:   imageURL,
		}, nil
	} else {
		return &Restaurant{}, errors.New("No Restaurants Found for " + location)
	}
}
