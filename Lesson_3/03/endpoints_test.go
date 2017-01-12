package main

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestEndpointes(t *testing.T) {
	fmt.Println("Running Endpoint Tester....")

	//Making a GET Request
	fmt.Println("Making a GET Request for /puppies...")
	url := "http://localhost:5000"

	resp, err := http.Get(url + "/puppies")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 1 FAILED: Could not make GET Request to web server")
		t.FailNow()
	} else {
		fmt.Println("Test 1 PASS: Succesfully Made GET Request to /puppies")
	}

	//#Making GET Requests to /puppies/id
	fmt.Println("Making GET requests to /puppies/id")

	id := 1
	for id <= 10 {
		resp, err = http.Get(url + "/puppies/" + strconv.Itoa(id))
		if err != nil {
			t.Fatal(err)
		}
		id++

		if resp.StatusCode != http.StatusOK {
			t.Log("Received an unsuccessful status code of", resp.StatusCode)
			t.Log("Test 2 FAILED: Could not make GET Request to /puppies/id")
			t.FailNow()
		}
	}

	fmt.Println("Test 2 PASS: Succesfully Made GET Request to /puppies/id")
}
