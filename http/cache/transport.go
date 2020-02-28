package cache

import (
	"log"
	"net/http"
)

func New(base http.RoundTripper, cache Cache) http.RoundTripper {
	// TODO: cache
	return &CacheTransport{base: base, cache: cache}
}

type CacheTransport struct {
	base  http.RoundTripper
	cache Cache
}

func (t *CacheTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Printf("RoundTrip: %#v", r)

	cache, err := t.cache.Get(r)
	if cache != nil {
		log.Printf("cache hit")
		return cache, nil
	}

	log.Printf("cache MISSED")
	resp, err := t.base.RoundTrip(r)
	if err != nil {
		return resp, err
	}
	if err := t.cache.Set(r, resp); err != nil {
		log.Printf("! failed to store cache: %s", err)
	}
	log.Printf("successfully cached")

	return resp, err
}
