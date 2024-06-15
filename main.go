package main

import (
	"io"
	"log"
	"net/http"

	"github.com/i7a7467/dev/db"
)

func main() {

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		res , err := db.GetCurrentTime()
		if err != nil {
			http.Error(w, "error occured.", http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, res + " server running.")
		}	
	}
	http.HandleFunc("/health", helloHandler)
	log.Println("server start at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
