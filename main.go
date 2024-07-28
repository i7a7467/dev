package main

import (
	"log"
	"net/http"

	"github.com/i7a7467/dev/handler"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthCheckHandler)
	mux.HandleFunc("/status", handler.StatusCheckHandler)
	mux.HandleFunc("/persons/{id}", handler.GetOnePersonHandler)

	log.Println("server start at port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
