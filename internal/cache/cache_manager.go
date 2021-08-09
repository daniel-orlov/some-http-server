package cache

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.Duration("cache_entry_ttl", 60*time.Second, "Cache entry TTL (def: 60 sec).")
	pflag.Int("ristretto_entry_cost", 1, "Ristretto entry cost per item saved in cache (def: 1).")
	pflag.Int("ristretto_num_counters", 100, "Ristretto number of keys to track for their frequency (def: 100).")
	pflag.Int("ristretto_buffer_items", 64, "Ristretto number of keys per Get buffer (def. 64).")
	pflag.Int("ristretto_max_cost", 100, "Ristretto maximum total cost of cache (def. 100).")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

type Store interface {
	Get(key interface{}) (interface{}, bool)
	Set(key interface{}, value interface{}) bool
	SetWithTTL(key, value interface{}) bool
}

type Options struct {
	RistrettoEntryCost   int64
	RistrettoNumCounters int64
	RistrettoMaxCost     int64
	RistrettoBufferItems int64
	TTL                  time.Duration
}

const (
	Ristretto = iota
)

func GetOptsFromEnv() *Options {
	if !viper.GetBool("cache_enable") {
		return nil
	}

	return &Options{
		TTL:                  viper.GetDuration("cache_entry_ttl"),
		RistrettoEntryCost:   viper.GetInt64("ristretto_entry_cost"),
		RistrettoNumCounters: viper.GetInt64("ristretto_num_counters"),
		RistrettoBufferItems: viper.GetInt64("ristretto_buffer_items"),
		RistrettoMaxCost:     viper.GetInt64("ristretto_max_cost"),
	}
}

// NewCacheStore creates a new instance of a cache store based on caching type.
func NewCacheStore(cachingType int, opts *Options) Store {
	if opts == nil {
		return nil
	}

	if cachingType == Ristretto {
		return NewRistrettoCache(opts)
	}

	return nil
}
