package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/i7a7467/dev/cache"
	"github.com/i7a7467/dev/handler"
	"github.com/valkey-io/valkey-go"
)

func main() {

    apiCache, err := cache.InitializeCache()
	if err != nil {
		log.Fatal(err)
	}

	err = cache.InitCache(0, 5*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

    handler := handler.NewHandler(apiCache)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthCheckHandler)
	mux.HandleFunc("/status", handler.StatusCheckHandler)
	mux.HandleFunc("/accounts/{id}", handler.GetOneAccountHandler)
	mux.HandleFunc("/accounts", handler.GetAccountsHandler) //bigcache
	mux.HandleFunc("/cache", handler.CacheTestHandler)
	mux.HandleFunc("/lruaccounts", handler.GetLruAccountsHandler) //golang-lru
	mux.HandleFunc("/cache/update", handler.CacheUpdateHandler) //for valkey sample

	// for valkey sample
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{os.Getenv("CACHE_DB_HOST") + ":" + os.Getenv("CACHE_DB_PORT")}, Username: os.Getenv("CACHE_DB_USER"), Password: os.Getenv("CACHE_DB_PASS")})

	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx := context.Background()
	err = client.Do(ctx, client.B().Set().Key("key").Value("valuedayo!!").Nx().Build()).Error()
	resp, err := client.Do(ctx, client.B().Get().Key("key").Build()).ToString();
    fmt.Println(resp)

	log.Println("server start at port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
