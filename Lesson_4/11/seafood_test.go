package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestEndpointes(t *testing.T) {
	fmt.Println("Running Endpoint Tester....")
	address := "http://localhost:5000"

	//GET AUTH CODE
	clientURL := address + "/clientOAuth"

	authCode := os.Getenv("code")
	if authCode == "" {
		t.Logf("Visit %s in your browser and get authorization code\n", clientURL)
		t.Log("Then write your test like that code=authorization-code go test")
		t.FailNow()
	}

	//TEST ONE GET TOKEN
	newContent := map[string]string{}
	var token string

	urlAddress := address + "/oauth/google"
	resp, err := http.PostForm(urlAddress, url.Values{"auth_code": {authCode}})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("Received an unsuccessful status code of", resp.StatusCode)
		goto ERROR
	}

	if err := json.NewDecoder(resp.Body).Decode(&newContent); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if newContent["token"] == "" {
		t.Error("No Token Received!")
		goto ERROR
	}

	token = newContent["token"]

ERROR:
	if t.Failed() {
		t.Fatal("Test 1 FAILED: Could not exchange auth code for a token")
	} else {
		fmt.Println("received token:", token)
		fmt.Println("test 1 PASS: Succesfully obtained token! ")
	}

}
