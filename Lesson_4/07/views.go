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
	route.HandleFunc("/token", verify(getAuthToken)).Methods("GET")

	log.Println("Serving HTTP on port", "5000")
	http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, route))

}

func getAuthToken(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(User)
	if !ok {
		http.Error(w, "Something wrong", 500)
		return
	}

	token, err := user.GenerateAuthToken(600)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		w.WriteHeader(200)
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
	user, ok := r.Context().Value("user").(User)
	if !ok {
		http.Error(w, "Something wrong", 500)
		return
	}

	Msg := fmt.Sprintf("Hello, %v!", user.Username)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"data": Msg}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func verify(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usernameOrToken, password, _ := r.BasicAuth()

		user := User{}
		//Try to see if it's a token first
		userID, err := user.VerifyAuthToken(usernameOrToken)
		if err != nil {
			DB.Where("username = ?", usernameOrToken).First(&user)
			if user.ID == 0 {
				http.Error(w, "User not found", 400)
				return
			} else if err := user.VerifyPassword(password); err != nil {
				http.Error(w, "Unable to verify password", 400)
				return
			}

		} else {
			DB.Where("id = ?", userID).First(&user)
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		fn(w, r)
	})
}
