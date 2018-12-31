package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// Test should be run on Prod VM.

func TestGetFromCache(t *testing.T) {
	res, err := getRESTDataFromCache("https://mock_url.com", getSecretLocal)
	assert.Nil(t, err)
	log.Printf("cached data: %v", res)
}

func TestGetSecret(t * testing.T)  {
	str, err := getSecret(testSecretUrl)
	assert.Nil(t, err)
	log.Printf("value: %v", str)
}
