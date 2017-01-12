package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type JSONPuppies struct {
	Puppies []JSONPuppy `json:"puppies"`
}

type JSONPuppy struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	route := mux.NewRouter().StrictSlash(true)

	route.HandleFunc("/", puppiesFunction).Methods("GET")
	route.HandleFunc("/puppies", puppiesFunction).Methods("GET", "POST")
	route.HandleFunc("/puppies/{id}", puppiesFunctionId).Methods("GET", "PUT", "DELETE")

	log.Println("Serving HTTP on port", "5000")
	http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, route))
}

func puppiesFunction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		//Call the method to Get all of the puppies
		w.Write(getAllPuppies())
	case "POST":
		//Call the method to make a new puppy
		fmt.Println("Making a New puppy")

		name := r.FormValue("name")
		description := r.FormValue("description")
		fmt.Println(name)
		fmt.Println(description)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(makeNewPuppy(name, description))
	}
}

func getAllPuppies() []byte {
	puppies := []Puppy{}
	DB.Find(&puppies)

	jsonPuppies := []JSONPuppy{}
	for _, puppy := range puppies {
		jsonPuppy := JSONPuppy{
			ID:          puppy.ID,
			Name:        puppy.Name,
			Description: puppy.Description,
		}
		jsonPuppies = append(jsonPuppies, jsonPuppy)
	}

	j, err := json.MarshalIndent(JSONPuppies{jsonPuppies}, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	return j
}

func makeNewPuppy(name, description string) []byte {
	puppy := Puppy{Name: name, Description: description}
	row := DB.Create(&puppy).Select("id").Row()

	var id uint
	row.Scan(&id)

	return getPuppy(int(id))
}

func getPuppy(id int) []byte {
	puppy := Puppy{}
	DB.First(&puppy, Puppy{ID: uint(id)})

	jsonPuppy := JSONPuppy{
		ID:          puppy.ID,
		Name:        puppy.Name,
		Description: puppy.Description,
	}

	j, err := json.MarshalIndent(jsonPuppy, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	return j
}

func puppiesFunctionId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	switch r.Method {
	case "GET":
		//Call the method to view a specific puppy
		w.Header().Set("Content-Type", "application/json")
		w.Write(getPuppy(id))

	case "PUT":
		//Call the method to edit a specific puppy
		name := r.FormValue("name")
		description := r.FormValue("description")
		updatePuppy(id, name, description)

	case "DELETE":
		//Call the method to remove a puppy
		deletePuppy(id)
	}
}

func updatePuppy(id int, name, description string) error {
	puppy := Puppy{}
	DB.First(&puppy, Puppy{ID: uint(id)})
	if puppy.ID == 0 {
		return errors.New("Not registered")
	}

	if name != "" {
		puppy.Name = name
	}

	if description != "" {
		puppy.Description = description
	}

	DB.Save(&puppy)

	fmt.Println("Updated a Puppy with id", id)

	return nil
}

func deletePuppy(id int) error {
	puppy := Puppy{}
	DB.First(&puppy, Puppy{ID: uint(id)})
	if puppy.ID == 0 {
		return errors.New("Not registered")
	}

	DB.Delete(&puppy)

	fmt.Println("Removing Puppy with id", id)

	return nil
}
