package main

import (
	"log"
	"net/http"

	"github.com/i7a7467/dev/handler"
)

func main() {

	http.HandleFunc("/health", handler.HealthCheckHandler)

	http.HandleFunc("/status", handler.StatusCheckHandler)
	
	log.Println("server start at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
