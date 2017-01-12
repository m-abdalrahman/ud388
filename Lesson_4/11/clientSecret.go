package main

import (
	"encoding/json"
	"log"
	"os"
)

var filePATH = "./client_secrets.json"

type CleintSecret struct {
	Web struct {
		Secret string `json:"client_secret"`
		ID     string `json:"client_id"`
	} `json:"Web"`
}

func GetCleintSecret() CleintSecret {
	file, err := os.Open(filePATH)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	clientSecret := CleintSecret{}
	if err := json.NewDecoder(file).Decode(&clientSecret); err != nil {
		log.Println(err)
	}

	return clientSecret
}
