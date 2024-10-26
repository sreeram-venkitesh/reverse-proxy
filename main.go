package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func handleRequest(rw http.ResponseWriter, r *http.Request) {
	targetHost := "http://localhost:9000"

	targetUrl := fmt.Sprintf("%s%s", targetHost, r.URL.Path)
	proxyRequest, err := http.NewRequest(r.Method, targetUrl, r.Body)
	if err != nil {
		http.Error(rw, "Error creating proxy", http.StatusInternalServerError)
		return
	}

	// Copying headers from the client's request to our proxied request
	for header, values := range r.Header {
		for _, value := range values {
			proxyRequest.Header.Add(header, value)
		}
	}

	// Proxy forwarding the request to target
	client := &http.Client{}
	res, err := client.Do(proxyRequest)
	if err != nil {
		http.Error(rw, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Copying headers from target server response to the proxy response
	for header, values := range res.Header {
		for _, value := range values {
			rw.Header().Set(header, value)
		}
	}

	rw.WriteHeader(res.StatusCode)

	io.Copy(rw, res.Body)
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
