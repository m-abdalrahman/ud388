package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type JSONUser struct {
	Username string `json:"username"`
}

func main() {
	route := mux.NewRouter().StrictSlash(true)

	route.HandleFunc("/clientOAuth", start).Methods("GET")
	route.HandleFunc("/oauth/{provider}", login).Methods("POST")
	route.HandleFunc("/token", verify(getAuthToken)).Methods("GET")
	route.HandleFunc("/users", newUser).Methods("POST")
	route.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	route.HandleFunc("/api/resource", verify(getResource)).Methods("GET")

	log.Println("Serving HTTP on port", "5000")
	http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, route))

}

func start(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/clientOAuth.html"))
	tmpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//STEP 1 - Parse the auth code
	authCode := r.FormValue("auth_code")
	fmt.Println("Step 1 - Complete, received auth code", authCode)

	if vars["provider"] == "google" {
		//STEP 2 - Exchange for a token
		oauthFlow := &oauth2.Config{
			ClientID:     GetCleintSecret().Web.ID,
			ClientSecret: GetCleintSecret().Web.Secret,
			RedirectURL:  "postmessage",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://accounts.google.com/o/oauth2/token",
			},
		}

		token, err := oauthFlow.Exchange(oauth2.NoContext, authCode)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			errMsg := fmt.Sprint("Failed to upgrade the authorization code.")
			if err := json.NewEncoder(w).Encode(map[string]string{"error": errMsg}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		//Check that the access token is valid.
		accessToken := token.AccessToken
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		fmt.Println("Step 2 Complete! Access Token : ", accessToken)

		//STEP 3 - Find User or make a new one

		//Get user info
		data := map[string]interface{}{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := data["name"].(string)
		picture := data["picture"].(string)
		email := data["email"].(string)

		//see if user exists, if it doesn't make a new one
		user := User{}

		DB.First(&user, User{Email: email})
		if user.Email == "" {
			user = User{Username: name, Picture: picture, Email: email}
			DB.Create(&user)
		}

		//STEP 4 - Make token
		authToken, err := user.GenerateAuthToken(600)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//STEP 5 - Send back token to the client
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"token": authToken}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Unrecoginized Provider", http.StatusUnauthorized)
	}
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

//middleware
func verify(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username_or_token, password, _ := r.BasicAuth()

		user := User{}
		//Try to see if it's a token first
		userID, err := user.VerifyAuthToken(username_or_token)
		if err != nil {
			DB.Where("username = ?", username_or_token).First(&user)
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
