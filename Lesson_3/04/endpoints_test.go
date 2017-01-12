package main

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestEndpointes(t *testing.T) {
	fmt.Println("Running Endpoint Tester....")
	url := "http://localhost:5000"
	clinet := &http.Client{}

	//Making a GET Request
	fmt.Println("Making a GET Request for /puppies...")

	resp, err := clinet.Get(url + "/puppies")
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

	//Making a POST Request
	fmt.Println("Making a POST Request to /puppies...")

	resp, err = clinet.Post(url+"/puppies", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 2 FAILED: Could not make POST Request to web server")
		t.FailNow()
	} else {
		fmt.Println("Test 2 PASS: Succesfully Made POST Request to /puppies")
	}

	//Making GET Requests to /puppies/id
	fmt.Println("Making GET requests to /puppies/id")

	for id := 1; id <= 10; id = id + 1 {
		resp, err = clinet.Get(url + "/puppies/" + strconv.Itoa(id))
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Log("Received an unsuccessful status code of", resp.StatusCode)
			t.Log("Test 3 FAILED: Could not make GET Requests to web server")
			t.FailNow()
		}

	}
	fmt.Println("Test 3 PASS: Succesfully Made GET Request to /puppies/id")

	//Making a PUT Request
	fmt.Println("Making PUT requests to /puppies/id")

	for id := 1; id <= 10; id = id + 1 {
		req, err := http.NewRequest("PUT", url+"/puppies/"+strconv.Itoa(id), nil)
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
		}

	}
	fmt.Println("Test 4 PASS: Succesfully Made PUT Request to /puppies/id")

	//Making a DELETE Request
	fmt.Println("Making DELETE requests to /puppies/id ...")

	for id := 1; id <= 10; id = id + 1 {
		req, err := http.NewRequest("DELETE", url+"/puppies/"+strconv.Itoa(id), nil)
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
		}

	}
	fmt.Println("Test 5 PASS: Succesfully Made DELETE Request to /puppies/id")

	fmt.Println("ALL TESTS PASSED!!")
}
