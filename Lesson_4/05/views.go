package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	route.HandleFunc("/users", newUser).Methods("POST")
	route.HandleFunc("/users/{id}", getUser).Methods("GET")
	route.HandleFunc("/resource", verify(getResource)).Methods("GET")
	route.HandleFunc("/bagels", verify(showAllBagels)).Methods("GET", "POST")

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
		fmt.Println("existing user")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "user already exists"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

func getResource(w http.ResponseWriter, r *http.Request) {
	uname := r.Context().Value("username")

	Msg := fmt.Sprintf("Hello, %s!", uname)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"data": Msg}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func showAllBagels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		bagels := []Bagel{}

		DB.Find(&bagels)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bagels); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// j, err := json.MarshalIndent(bagels, "", "  ")
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// w.Write(j)

	case "POST":
		bagel := Bagel{}
		if err := json.NewDecoder(r.Body).Decode(&bagel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		DB.Create(&bagel)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		if err := json.NewEncoder(w).Encode(bagel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func verify(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()

		fmt.Println("Looking for user", username)

		user := User{}
		DB.Where("username = ?", username).First(&user)
		if user.ID == 0 {
			http.Error(w, "User not found", 400)
			return
		} else if err := user.VerifyPassword(password); err != nil {
			http.Error(w, "Unable to verify password", 400)
			return
		}

		ctx := context.WithValue(r.Context(), "username", user.Username)
		r = r.WithContext(ctx)

		fn(w, r)
	})
}
