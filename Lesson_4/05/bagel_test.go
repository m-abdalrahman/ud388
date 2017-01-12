package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestEndpointes(t *testing.T) {
	fmt.Println("Running Endpoint Tester....")
	address := "http://localhost:5000"

	//TEST 1 TRY TO MAKE A NEW USER
	newUser := User{Username: "TinnyTim", PasswordHash: "Udacity"}

	data := new(bytes.Buffer)

	if err := json.NewEncoder(data).Encode(newUser); err != nil {
		log.Println(err)
	}

	urlAdd := address + "/users"
	resp, err := http.Post(urlAdd, "application/json", data)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 1 FAILED: Could not make a new user")
		t.FailNow()
	} else {
		fmt.Println("Test 1 PASS: Succesfully made a new user")
	}

	//TEST 2 ADD NEW BAGELS TO THE DATABASE
	newBagel := Bagel{
		Name:        "plain",
		Picture:     "http://bonacbagel.weebly.com/uploads/4/0/5/4/40548977/s318635836612132814_p1_i1_w240.jpeg",
		Description: "Old-Fashioned Plain Bagel",
		Price:       "$1.99",
	}
	if err := json.NewEncoder(data).Encode(newBagel); err != nil {
		log.Println(err)
	}

	urlAdd = address + "/bagels"
	req, err := http.NewRequest("POST", urlAdd, data)
	if err != nil {
		log.Println(err)
	}

	req.SetBasicAuth("TinnyTim", "Udacity")
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 2 FAILED: Could not add new bagels")
		t.FailNow()
	} else {
		fmt.Println("Test 2 PASS: Succesfully made new bagels")
	}

	//TEST 3 TRY TO READ BAGELS WITH INVALID CREDENTIALS
	req, err = http.NewRequest("GET", urlAdd, nil)
	if err != nil {
		log.Println(err)
	}

	req.SetBasicAuth("TinnyTim", "youdacity")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Log("Security Flaw: able to log in with invalid credentials")
		t.Log("Test 3 FAILED")
		t.FailNow()
	} else {
		fmt.Println("Test 3 PASS: App checks against invalid credentials")
	}

	//TEST 4 TRY TO READ BAGELS WITH VALID CREDENTIALS
	req, err = http.NewRequest("GET", urlAdd, nil)
	if err != nil {
		log.Println(err)
	}

	req.SetBasicAuth("TinnyTim", "Udacity")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Log("Test 4 FAILED")
		t.FailNow()
	} else {
		fmt.Println("Test 4 PASS: Logged in User can view /bagels")
	}

}
