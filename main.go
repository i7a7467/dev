package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx := context.Background()
	// SET key val NX
	err = client.Do(ctx, client.B().Set().Key("key2").Value("valuedayo!!").Nx().Build()).Error()
	// HGETALL hm
	//hm, err := client.Do(ctx, client.B().Hgetall().Key("hm").Build()).AsStrMap()

	resp, err := client.Do(ctx, client.B().Get().Key("key2").Build()).ToString();
    fmt.Println(resp)


	log.Println("server start at port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
