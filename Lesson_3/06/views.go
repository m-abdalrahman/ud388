package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter().StrictSlash(true)

	route.HandleFunc("/restaurants", allRestaurantsHandler).Methods("GET", "POST")
	route.HandleFunc("/restaurants/{id}", restaurantHandler).Methods("GET", "PUT", "DELETE")

	log.Println("Serving HTTP on port", "5000")
	http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, route))

}

func allRestaurantsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		//RETURN ALL RESTAURANTS IN DATABASE
		restaurant := []Restaurant{}
		DB.Find(&restaurant)

		j, err := json.MarshalIndent(restaurant, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)

	case "POST":
		location := r.FormValue("location")
		mealType := r.FormValue("mealType")

		w.Header().Set("Content-Type", "application/json")

		restaurantInfo, err := FindARestaurant(mealType, location)
		if err != nil {
			msgErr := fmt.Sprintf("No Restaurants Found for %s in %s", mealType, location)

			jErr := struct {
				Error string `json:"error"`
			}{msgErr}

			j, err := json.MarshalIndent(jErr, "", "  ")
			if err != nil {
				fmt.Println(err)
			}

			w.WriteHeader(404)
			w.Write(j)
		} else {
			//MAKE A NEW RESTAURANT AND STORE IT IN DATABASE
			DB.Create(&restaurantInfo)

			j, err := json.MarshalIndent(restaurantInfo, "", "  ")
			if err != nil {
				fmt.Println(err)
			}

			w.WriteHeader(201)
			w.Write(j)

		}
	}
}

func restaurantHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	restaurant := Restaurant{}
	DB.First(&restaurant, Restaurant{ID: uint(id)})

	switch r.Method {
	case "GET":
		//RETURN A SPECIFIC RESTAURANT
		j, err := json.MarshalIndent(restaurant, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)

	case "PUT":
		//UPDATE A SPECIFIC RESTAURANT
		address := r.FormValue("address")
		image := r.FormValue("image")
		name := r.FormValue("name")

		if address != "" {
			restaurant.Address = address
		}
		if image != "" {
			restaurant.Image = image
		}
		if name != "" {
			restaurant.Name = name
		}

		DB.Save(&restaurant)

		j, err := json.MarshalIndent(restaurant, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)

	case "DELETE":
		//DELETE A SPECFIC RESTAURANT
		DB.Delete(&restaurant)

		w.Write([]byte("Restaurant Deleted"))
	}

}
