package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
    "github.com/hashicorp/golang-lru/v2/expirable"
    "github.com/i7a7467/dev/model"
	"github.com/dgraph-io/ristretto"

)

var cache *bigcache.BigCache

var lruCache *expirable.LRU[string, interface{}]
type Cache struct {
    bigCache *bigcache.BigCache
}
type LruAccountsCache struct {
    lruCache *expirable.LRU[string, model.Accounts]
}

type Rcache struct {
    ristrettoCache *ristretto.Cache
}

func InitializeCache() (*Cache ,error) {
	config := bigcache.Config{
        LifeWindow: 		   60 * time.Second,
        CleanWindow:           1 * time.Minute,
        MaxEntriesInWindow:    1000 * 10 * 60,
        MaxEntrySize:          500,
        HardMaxCacheSize:      1024,
        OnRemove:              nil,
        OnRemoveWithReason:    nil,
        Logger:                nil,
		Shards: 			   2,
    }
    var err error
    cache, err = bigcache.New(context.Background(), config)
    if err != nil {
        return nil ,err
    }
    return &Cache{bigCache: cache} ,nil
}

func InitCache(size int, ttl time.Duration) error {
    var err error
    lruCache = expirable.NewLRU[string, interface{}](size, nil, ttl)
    return err
}

func (c *Cache) SetCache(key string, value []byte) error {
    return c.bigCache.Set(key, value)
}

func (c *Cache) GetCache(key string) ([]byte, error) {
    return c.bigCache.Get(key)
}

func (c *Cache) IsInitialized() bool {
    return c.bigCache != nil
}

// Add adds an entry to the cache
func Add(key string, value interface{}) {
    lruCache.Add(key, value)
}

// Get retrieves an entry from the cache
func Get(key string) (interface{}, bool) {
    return lruCache.Get(key)
}

