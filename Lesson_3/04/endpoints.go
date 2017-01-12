package main

import (
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

	route.HandleFunc("/puppies", puppiesFunction).Methods("GET", "POST")
	route.HandleFunc("/puppies/{id}", puppiesFunctionId).Methods("GET", "PUT", "DELETE")

	log.Println("Serving HTTP on port", "5000")
	http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, route))
}

func puppiesFunction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintln(w, getAllPuppies())
	case "POST":
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, makeANewPuppy())
	}
}

func puppiesFunctionId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	switch r.Method {
	case "GET":
		fmt.Fprintln(w, getPuppy(id))
	case "PUT":
		fmt.Fprintln(w, updatePuppy(id))
	case "DELETE":
		fmt.Fprintln(w, deletePuppy(id))
	}
}

func getAllPuppies() string {
	return "Getting All the puppies!"
}

func makeANewPuppy() string {
	return "Creating A New Puppy!"
}

func getPuppy(id int) string {
	return fmt.Sprintln("Getting Puppy with id", id)
}

func updatePuppy(id int) string {
	return fmt.Sprintln("Updating a Puppy with id", id)
}

func deletePuppy(id int) string {
	return fmt.Sprintln("Removing Puppy with id", id)
}
