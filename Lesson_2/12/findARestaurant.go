package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	foursquareClientID     = ""
	foursquareClientSecret = ""
)

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

func main() {
	FindARestaurant("Pizza", "Tokyo, Japan")
	FindARestaurant("Tacos", "Jakarta, Indonesia")
	FindARestaurant("Tapas", "Maputo, Mozambique")
	FindARestaurant("Falafel", "Cairo, Egypt")
	FindARestaurant("Spaghetti", "New Delhi, India")
	FindARestaurant("Cappuccino", "Geneva, Switzerland")
	FindARestaurant("Sushi", "Los Angeles, California")
	FindARestaurant("Steak", "La Paz, Bolivia")
	FindARestaurant("Gyros", "Sydney Australia")
}

func FindARestaurant(mealType, location string) {
	//1. Use getGeocodeLocation to get the latitude and longitude coordinates of the location string.
	lat, lng := GetGeocodeLocation(location)

	//2.  Use foursquare API to find a nearby restaurant with the latitude, longitude, and mealType strings.
	//HINT: format for url will be something like https://api.foursquare.com/v2/venues/search?client_id=CLIENT_ID&client_secret=CLIENT_SECRET&v=20130815&ll=40.7,-74&query=sushi
	urlAddress := fmt.Sprintf("https://api.foursquare.com/v2/venues/search?client_id=%s&client_secret=%s&v=20130815&ll=%s,%s&query=%s",
		foursquareClientID, foursquareClientSecret,
		strconv.FormatFloat(lat, 'f', -1, 64), strconv.FormatFloat(lng, 'f', -1, 64), mealType)

	resp, err := http.Get(urlAddress)
	if err != nil {
		log.Println("Response:", err)
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

		// restaurant address
		restaurantAddress := strings.Join(result.Response.Venues[0].Location.FormattedAddress, " ")

		//4.  Get a  300x300 picture of the restaurant using the venue_id (you can change this by altering the 300x300 value in the URL or replacing it with 'orginal' to get the original picture
		imgURL := fmt.Sprintf("https://api.foursquare.com/v2/venues/%s/photos?client_id=%s&v=20150603&client_secret=%s",
			result.Response.Venues[0].ID, foursquareClientID, foursquareClientSecret)

		resp, err := http.Get(imgURL)
		if err != nil {
			log.Println("Response:", err)
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

		fmt.Println("Restaurant Name:", restaurantName)
		fmt.Println("Restaurant Address:", restaurantAddress)
		fmt.Println("Image:", imageURL)
		fmt.Println()
	} else {
		fmt.Println("No Restaurants Found for", location)
		fmt.Println()
	}
}
