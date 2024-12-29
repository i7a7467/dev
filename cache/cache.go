package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dgraph-io/ristretto"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/i7a7467/dev/model"
	"github.com/valkey-io/valkey-go"
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

type ValkeyConfig struct {
    InitAddress []string
    Username    string
    Password    string
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

func SetValkey(client valkey.Client,key string, value []model.Person) {
    data, err := json.Marshal(value)
    if err != nil {
        // handle error
        return
    }
    resp := client.Do(context.Background(), client.B().Set().Key(key).Value(string(data)).Nx().Build())
    if resp.Error() != nil {
        // handle error
        fmt.Println(resp.Error())
        return
    }
}

func GetValkey(client valkey.Client, key string) ([]model.Person, error) {
    result,err := client.Do(context.Background(), client.B().Get().Key(key).Build()).ToString()

    var accounts []model.Person
    if err != nil {
        // handle error
        return accounts, err
    }
    err = json.Unmarshal([]byte(result), &accounts)
    if err != nil {
        return accounts, err
    }
    return accounts, nil
}

func CacheDBConn() (valkey.Client, error) {
    // for valkey sample
    cacheDBHosts :=  os.Getenv("CACHE_DB_HOST")
    sliceDBList := strings.Split(cacheDBHosts, ",")
	for i := range sliceDBList {
		sliceDBList[i] = strings.TrimSpace(sliceDBList[i]) + ":" + os.Getenv("CACHE_DB_PORT")
	}
    client, err := valkey.NewClient(valkey.ClientOption{InitAddress: sliceDBList, Username: os.Getenv("CACHE_DB_USER"), Password: os.Getenv("CACHE_DB_PASS")})

    if err != nil {
        //panic(err) //return で返す
        return nil, err
    }
    defer client.Close()

    return client, nil
    
}