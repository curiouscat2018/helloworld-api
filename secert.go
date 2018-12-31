package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getSecret(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	accessToken, err := getKeyVaultAccessToken()
	if err != nil {
		return "", err
	}

	client := http.Client{}
	req.Header.Set("Authorization", "Bearer " + accessToken)
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		rawBody, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("invalid http status: %v body: %v", response.StatusCode, string(rawBody))
	}

	responseJSON := struct {
		Value string
	}{}

	if err := json.NewDecoder(response.Body).Decode(&responseJSON); err != nil {
		return "", err
	}

	log.Println("successfully get secret")
	return responseJSON.Value, nil
}

func getKeyVaultAccessToken() (string, error) {
	req, err := http.NewRequest("GET", "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fvault.azure.net", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Metadata", "true")
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		rawBody, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("invalid http status: %v body: %v", response.Status, string(rawBody))
	}

	responseJSON := struct {
		Access_token string
	}{}

	if err := json.NewDecoder(response.Body).Decode(&responseJSON); err != nil {
		return "", err
	}

	log.Println("successfully get access_token")
	return responseJSON.Access_token, nil
}

