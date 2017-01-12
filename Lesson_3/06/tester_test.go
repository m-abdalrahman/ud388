package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func TestEndpointes(t *testing.T) {
	fmt.Println("Running Endpoint Tester....")
	urlAdd := "http://localhost:5000"

	clinet := &http.Client{}

	//TEST ONE -- CREATE NEW RESTAURANTS
	fmt.Println("Test 1: Creating new Restaurants......")

	resp, err := clinet.PostForm(urlAdd+"/restaurants", url.Values{"location": {"Buenos Aires Argentina"}, "mealType": {"Sushi"}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error("Received an unsuccessful status code of", resp.StatusCode)
		goto ERROR
	}

	fmt.Println(string(body))

	//
	resp, err = clinet.PostForm(urlAdd+"/restaurants", url.Values{"location": {"Denver Colorado"}, "mealType": {"Soup"}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error("Received an unsuccessful status code of", resp.StatusCode)
		goto ERROR
	}

	fmt.Println(string(body))

	//
	resp, err = clinet.PostForm(urlAdd+"/restaurants", url.Values{"location": {"Prague Czech Republic"}, "mealType": {"Crepes"}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error("Received an unsuccessful status code of", resp.StatusCode)
		goto ERROR
	}

	fmt.Println(string(body))

	//
	resp, err = clinet.PostForm(urlAdd+"/restaurants", url.Values{"location": {"Shanghai China"}, "mealType": {"Sandwiches"}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error("Received an unsuccessful status code of", resp.StatusCode)
		goto ERROR
	}

	fmt.Println(string(body))

	//
	resp, err = clinet.PostForm(urlAdd+"/restaurants", url.Values{"location": {"Nairobi Kenya"}, "mealType": {"Pizza"}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error("Received an unsuccessful status code of", resp.StatusCode)
		goto ERROR
	}

	fmt.Println(string(body))

ERROR:
	if t.Failed() {
		t.Fatal("Test 1 FAILED: Could not add new restaurants")
	} else {
		fmt.Println("Test 1 PASS: Succesfully Made all new restaurants")
	}

	//TEST TWO -- READ ALL RESTAURANTS
	fmt.Println("Attempting Test 2: Reading all Restaurants...")

	resp, err = clinet.Get(urlAdd + "/restaurants")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(body))

	restaurants := []Restaurant{}
	err = json.Unmarshal(body, &restaurants)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 2 FAILED: Could not retrieve restaurants from server")
		t.FailNow()
	} else {
		fmt.Println("Test 2 PASS: Succesfully read all restaurants")
	}

	//TEST THREE -- READ A SPECIFIC RESTAURANT
	fmt.Println("Attempting Test 3: Reading the last created restaurant...")

	restaurantID := restaurants[len(restaurants)-1].ID

	resp, err = clinet.Get(urlAdd + "/restaurants/" + strconv.Itoa(int(restaurantID)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 3 FAILED: Could not retrieve restaurant from server")
		t.FailNow()
	} else {
		fmt.Println("Test 3 PASS: Succesfully read last restaurant")
	}

	//TEST FOUR -- UPDATE A SPECIFIC RESTAURANT
	fmt.Println("Attempting Test 4: Changing the name, image, and address of the first restaurant to Udacity...")

	restaurantID = restaurants[0].ID
	req, err := http.NewRequest("PUT",
		urlAdd+"/restaurants/"+strconv.Itoa(int(restaurantID))+"?name=Udacity&address=2465+Latham+Street+Mountain+View+CA&image=https://media.glassdoor.com/l/70/82/fc/e8/students-first.jpg", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = clinet.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 4 FAILED: Could not update restaurant from server")
		t.FailNow()
	} else {
		fmt.Println("Test 4 PASS: Succesfully updated first restaurant")
	}

	//TEST FIVE -- DELETE SECOND RESTARUANT
	fmt.Println("Attempting Test 5: Deleteing the second restaurant from the server...")
	restaurantID = restaurants[1].ID
	req, err = http.NewRequest("DELETE",
		urlAdd+"/restaurants/"+strconv.Itoa(int(restaurantID)), nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = clinet.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		t.Log("Received an unsuccessful status code of", resp.StatusCode)
		t.Log("Test 5 FAILED: Could not delete restaurant from server")
		t.FailNow()
	} else {
		fmt.Println("Test 5 PASS: Succesfully delete second restaurant")
	}

	fmt.Println("ALL TESTS PASSED!")
}
