package cache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/cache"
)

const CacheGCSec = 3600
const CachePersistTime = time.Hour

type Cache struct {
	gcSec       int
	persistTime time.Duration
	cache       cache.Cache
}

func NewCache(GCSec int, persistTime time.Duration) (*Cache, error) {
	c, err := cache.NewCache("memory", `{"interval":`+strconv.Itoa(GCSec)+`}`)
	if err != nil {
		return nil, err
	}
	wrapper := &Cache{}
	wrapper.gcSec = GCSec
	wrapper.persistTime = persistTime
	wrapper.cache = c

	return wrapper, nil
}

type GetRestDataFunctor func(string) (string, error)

func (c *Cache) GetRESTDataFromCache(url string, functor GetRestDataFunctor) (data string, isCacheFound bool, err error) {
	if c.cache.IsExist(url) {
		res, ok := c.cache.Get(url).(string)
		if !ok {
			return "", true, fmt.Errorf("not able to cast value to string")
		}
		return res, true, nil
	}

	secret, err := functor(url)
	if err != nil {
		return "", false, err
	}

	if err := c.cache.Put(url, secret, c.persistTime); err != nil {
		return "", false, nil
	}

	return secret, false, nil
}
