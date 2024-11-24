package main

import (
	"log"
	"net/http"
	"time"

	"github.com/i7a7467/dev/cache"
	"github.com/i7a7467/dev/handler"
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

	log.Println("server start at port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
