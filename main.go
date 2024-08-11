package main

import (
	"log"
	"net/http"

	"github.com/i7a7467/dev/handler"
	"github.com/i7a7467/dev/cache"

)

func main() {

    apiCache, err := cache.InitializeCache()
	if err != nil {
		log.Fatal(err)
	}

    handler := handler.NewHandler(apiCache)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthCheckHandler)
	mux.HandleFunc("/status", handler.StatusCheckHandler)
	mux.HandleFunc("/accounts/{id}", handler.GetOneAccountHandler)
	mux.HandleFunc("/accounts", handler.GetAccountsHandler) //cache
	mux.HandleFunc("/cache", handler.CacheTestHandler)

	log.Println("server start at port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
