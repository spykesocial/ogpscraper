package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yashdiniz/ogpscraper/internal/metaparser"
)

const (
	defaultTtl = 24
)

var (
	cache = redis.NewClient(&redis.Options{})
)

func IsInitialized() bool {
	return cache != nil
}

// GetInstance returns a client if initialized otherwise throws an error
func GetInstance() (*redis.Client, error) {
	if !IsInitialized() {
		return nil, errors.New("cannot get instance of the cache - not intialized")
	}
	return cache, nil
}

// GetPageId gets the key used to store URL page info in cache
func GetPageId(u *url.URL) string {
	return fmt.Sprintf("%v%v", u.Host, u.Path)
}

/*
Checks the cache given a URL, if there is an entry it unmarshalls the JSON and returns
a pointer to the data
*/
func CheckCacheForPage(address *url.URL) ([]metaparser.MetaTag, error) {
	if cache == nil {
		return nil, errors.New("cache not initialized")
	}
	key := GetPageId(address)
	result, e := cache.Get(cache.Context(), key).Result()
	if e != nil {
		log.Printf("cache miss for %s\n", key)
		return nil, nil
	}
	log.Printf("cache hit for %s\n", key)
	tags := []metaparser.MetaTag{}
	marshallError := json.Unmarshal([]byte(result), &tags)
	if marshallError != nil {
		return nil, e
	}
	return tags, nil
}

/*
function to cache the page metadata in the redis cache.
// TODO: add configurable TTL for cache life
*/
func SetCachePageMetaData(tags []metaparser.MetaTag, address *url.URL, ttl int64) error {
	if cache == nil {
		return errors.New("cache not initialized")
	}

	pageId := GetPageId(address)
	log.Printf("caching meta data for %s\n", pageId)
	if serialized, marshalError := json.Marshal(tags); marshalError == nil {
		ctx := cache.Context()
		if cacheError := cache.Set(ctx, pageId, string(serialized), time.Duration(ttl)*time.Hour).Err(); cacheError == nil {
			log.Printf("cached metadata for page %s\n", pageId)
			return nil
		} else {
			return cacheError
		}
	} else {
		log.Println("could not marshal tags into object")
		return marshalError
	}
}
