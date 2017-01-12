package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func TestEndpointes(t *testing.T) {
	fmt.Println("Running Endpoint Tester....")
	urlAdd := "http://localhost:5000"

	clinet := &http.Client{}

	//Making a POST Request
	fmt.Println("Making a POST Request to /puppies...")

	resp, err := clinet.PostForm(urlAdd+"/puppies",
		url.Values{"name": {"Fido"}, "description": {"Playful Little Puppy"}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var puppy JSONPuppy
	err = json.NewDecoder(resp.Body).Decode(&puppy)
	if err != nil {
		t.Fatal(err)
	}

	//get last id
	puppyID := puppy.ID

	if resp.StatusCode != http.StatusCreated {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 1 FAILED: Could not make POST Request to web server")
		t.FailNow()
	} else {
		fmt.Println("Test 1 PASS: Succesfully Made POST Request to /puppies")
	}

	//Making a GET Request
	fmt.Println("Making a GET Request to /puppies...")

	resp, err = clinet.Get(urlAdd + "/puppies")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 2 FAILED: Could not make GET Request to web server")
		t.FailNow()
	} else {
		fmt.Println("Test 2 PASS: Succesfully Made GET Request to /puppies")
	}

	//Making GET Requests to /puppies/id
	fmt.Println("Making GET requests to /puppies/id")

	resp, err = clinet.Get(urlAdd + "/puppies/" + strconv.Itoa(int(puppyID)))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 3 FAILED: Could not make GET Requests to web server")
		t.FailNow()
	} else {
		fmt.Println("Test 3 PASS: Succesfully Made GET Request to /puppies/id")
	}

	//Making a PUT Request
	fmt.Println("Making PUT requests to /puppies/id")

	req, err := http.NewRequest("PUT",
		urlAdd+"/puppies/"+strconv.Itoa(int(puppyID))+"?name=wilma&description=A+sleepy+bundle+of+joy", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = clinet.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 4 FAILED: Could not make PUT Request to web server")
		t.FailNow()
	} else {
		fmt.Println("Test 4 PASS: Succesfully Made PUT Request to /puppies/id")
	}

	//Making a DELETE Request
	fmt.Println("Making DELETE requests to /puppies/id ...")

	req, err = http.NewRequest("DELETE",
		urlAdd+"/puppies/"+strconv.Itoa(int(puppyID)), nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = clinet.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 5 FAILED: Could not make DELETE Requests to web server")
		t.FailNow()
	} else {
		fmt.Println("Test 5 PASS: Succesfully Made DELETE Request to /puppies/id")
	}

	fmt.Println("ALL TESTS PASSED!!")
}
