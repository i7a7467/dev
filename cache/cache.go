package cache

import (
	"context"
	// "encoding/json"
	// "errors"
	// "fmt"
	"time"

	"github.com/allegro/bigcache/v3"
	// "github.com/i7a7467/dev/model"
)

var cache *bigcache.BigCache

type Cache struct {
    bigCache *bigcache.BigCache
}

func InitializeCache() (*Cache ,error) {
	config := bigcache.Config{
        LifeWindow: 		   1 * time.Second,
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

func (c *Cache) SetCache(key string, value []byte) error {
    return c.bigCache.Set(key, value)
}

func (c *Cache) GetCache(key string) ([]byte, error) {
    return c.bigCache.Get(key)
}

func (c *Cache) IsInitialized() bool {
    return c.bigCache != nil
}
