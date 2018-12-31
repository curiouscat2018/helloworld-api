package main

import (
	"fmt"
	"log"
	"time"
)

const cacheGCSec = 3600
const cachePersistTime = time.Hour

func getRESTDataFromCache(url string, functor func(string) (string, error)) (string, error) {
	if inmemoryCache.IsExist(url) {
		res, ok := inmemoryCache.Get(url).(string)
		if !ok {
			return "", fmt.Errorf("not able to cast value to string")
		}
		log.Printf("cache exist. successfully get from cache. url: %v", url)
		return res, nil
	}

	secert, err := functor(url)
	if err != nil {
		return "", err
	}

	if err := inmemoryCache.Put(url, secert, cachePersistTime); err != nil {
		return "", nil
	}

	log.Printf("cache not exist. successfully get and put into cache: %v", url)
	return secert, nil
}
