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

	//TEST 1: TRY TO REGISTER A NEW USER
	newUser := User{Username: "Peter", PasswordHash: "Pan"}

	data := new(bytes.Buffer)

	if err := json.NewEncoder(data).Encode(newUser); err != nil {
		log.Println(err)
	}

	urlAddress := address + "/users"
	resp, err := http.Post(urlAddress, "application/json", data)
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

	//TEST 2: OBTAIN A TOKEN
	newContent := map[string]string{}
	var token string

	urlAddress = address + "/token"
	req, err := http.NewRequest("GET", urlAddress, nil)
	if err != nil {
		log.Println(err)
	}

	req.SetBasicAuth("Peter", "Pan")
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Error("Received an unsuccessful status code of", resp.StatusCode)
		goto ERROR
	}

	if err := json.NewDecoder(resp.Body).Decode(&newContent); err != nil {
		t.Error(err)
		goto ERROR
	}
	token = newContent["token"]

ERROR:
	if t.Failed() {
		t.Fatal("Test 2 FAILED: Could not exchange user credentials for a token")
	} else {
		fmt.Println("received token:", token)
		fmt.Println("Test 2 PASS: Succesfully obtained token!")
	}

	//TEST 3: TRY TO ADD PRODUCS TO DATABASE
	newItem := Product{Name: "apple", Category: "fruit", Price: "$.99"}
	if err := json.NewEncoder(data).Encode(newItem); err != nil {
		log.Println(err)
	}

	urlAddress = address + "/products"

	req, err = http.NewRequest("POST", urlAddress, data)
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth(token, "blank")
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 3 FAILED: Could not add new products")
		t.FailNow()
	} else {
		fmt.Println("Test 3 PASS: Succesfully added new products")
	}

	//TEST 4: TRY ACCESSING ENDPOINT WITH AN INVALID TOKEN
	req, err = http.NewRequest("GET", urlAddress, nil)
	if err != nil {
		log.Println(err)
	}

	req.SetBasicAuth(token+"malpractice", "blank")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 4 FAILED: able to log in with invalid token")
		t.FailNow()
	} else {
		fmt.Println("Test 4 PASS: App checks against invalid credentials")
	}

	//TEST 5: TRY TO VIEW ALL PRODUCTS IN DATABASE
	products := []Product{} //for Decode from resp.Body

	var countProducts int
	DB.Table("product").Count(&countProducts)

	req, err = http.NewRequest("GET", urlAddress, nil)
	if err != nil {
		log.Println(err)
	}

	req.SetBasicAuth(token, "blank")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Error("Received an unsuccessful status code of", resp.StatusCode)
		goto ERROR2
	}

	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		log.Println(err)
	}

	//check if number of products from request ==
	//number of products inside database
	if len(products) != countProducts {
		t.Error("number of products from request not equal number of products inside database")
		goto ERROR2
	}

ERROR2:
	if t.Failed() {
		t.Fatal("Test 5 FAILED:")
	} else {
		fmt.Println("Test 5 PASS: viewed all products")
	}

	//TEST 6: TRY TO VIEW A SPECIFIC CATEGORY OF PRODUCTS
	var countFruits int
	DB.Model(&Product{}).Where("category = ?", "fruit").Count(&countFruits)

	urlAddress = address + "/products/fruit"
	req, err = http.NewRequest("GET", urlAddress, nil)
	if err != nil {
		log.Println(err)
	}

	req.SetBasicAuth(token, "blank")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Error("Received an unsuccessful status code of", resp.StatusCode)
		goto ERROR3
	}

	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		log.Println(err)
	}

	//check if number of products by category from request ==
	//number of products by category inside database
	if len(products) != countFruits {
		t.Error("number of products from request not equal number of products inside database")
		goto ERROR2
	}

ERROR3:
	if t.Failed() {
		t.Fatal("Test 6 FAILED: ")
	} else {
		fmt.Println("Test 6 PASS: viewed all categories")
	}
}
