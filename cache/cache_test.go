package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCache_GetRESTDataFromCache(t *testing.T) {
	inmemoryCache, err := NewCache(60, 1*time.Second)
	assert.Nil(t, err)
	testData := "testdata123"

	// First time
	res, found, err := inmemoryCache.GetRESTDataFromCache("https://mock_url.com",
		func(s2 string) (s string, e error) {
			return testData, nil
		})

	assert.Nil(t, err)
	assert.False(t, found)
	assert.Equal(t, res, testData)

	// Cache exist
	res, found, err = inmemoryCache.GetRESTDataFromCache("https://mock_url.com",
		func(s2 string) (s string, e error) {
			return testData, nil
		})
	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, res, testData)

	// Cache not exist
	time.Sleep(time.Second)
	res, found, err = inmemoryCache.GetRESTDataFromCache("https://mock_url.com",
		func(s2 string) (s string, e error) {
			return testData, nil
		})
	assert.Nil(t, err)
	assert.False(t, found)
	assert.Equal(t, res, testData)
}
