package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type JSONUser struct {
	Username string `json:"username"`
}

func main() {
	route := mux.NewRouter().StrictSlash(true)

	route.HandleFunc("/api/users", newUser).Methods("POST")
	route.HandleFunc("/api/users/{id}", getUser).Methods("GET")

	log.Println("Serving HTTP on port", "5000")
	http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, route))

}

func newUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Username == "" || user.PasswordHash == "" {
		w.WriteHeader(400) // missing arguments
		return
	}

	userQuery := User{}
	DB.First(&userQuery, User{Username: user.Username})
	if userQuery.Username != "" {
		w.WriteHeader(400) // existing user
		return
	}

	newUser := User{Username: user.Username}
	err := newUser.HashPassword(user.PasswordHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	DB.Create(&newUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	if err = json.NewEncoder(w).Encode(JSONUser{Username: user.Username}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	user := User{}
	DB.First(&user, User{ID: uint(id)})
	if user.ID == 0 {
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(JSONUser{Username: user.Username}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
