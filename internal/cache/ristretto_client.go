package cache

import (
	"github.com/dgraph-io/ristretto"
)

type RistrettoClient struct {
	opts  *Options
	cache *ristretto.Cache
}

// NewRistrettoCache creates a new cache instance.
func NewRistrettoCache(opts *Options) *RistrettoClient {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: opts.RistrettoNumCounters,
		MaxCost:     opts.RistrettoMaxCost,
		BufferItems: opts.RistrettoBufferItems,
	})
	if err != nil {
		panic(err)
	}

	return &RistrettoClient{cache: cache, opts: opts}
}

func (rc *RistrettoClient) Get(key interface{}) (interface{}, bool) {
	return rc.cache.Get(key)
}

func (rc *RistrettoClient) Set(key, value interface{}) bool {
	return rc.cache.Set(key, value, rc.opts.RistrettoEntryCost)
}

func (rc *RistrettoClient) SetWithTTL(key, value interface{}) bool {
	return rc.cache.SetWithTTL(key, value, rc.opts.RistrettoEntryCost, rc.opts.TTL)
}
