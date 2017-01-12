package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter().StrictSlash(true)

	route.HandleFunc("/puppies", puppiesFunction)
	route.HandleFunc("/puppies/{id}", puppiesFunctionId)

	log.Println("Serving HTTP on port", "5000")
	http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, route))
}

func puppiesFunction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Yes, puppies!")
}

func puppiesFunctionId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, "This method will act on the puppy with id", vars["id"])
}
