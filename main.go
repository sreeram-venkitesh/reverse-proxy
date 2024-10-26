package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	fmt.Fprintf(w, "Hello world!")
}

func main() {
	log.Printf("Reverse proxy is starting...")

	port := 8080

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        http.HandlerFunc(handleRequest),
	}
	
	log.Printf("Started server on port: %d\n", port)
	log.Fatal(s.ListenAndServe())
}
