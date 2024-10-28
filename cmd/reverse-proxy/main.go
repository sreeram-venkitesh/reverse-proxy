package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sreeram-venkitesh/reverse-proxy/pkg/config"
	"github.com/sreeram-venkitesh/reverse-proxy/pkg/proxy"
)

func main() {
	log.Printf("Reverse proxy is starting...")

	runtimeConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", runtimeConfig.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        http.HandlerFunc(proxy.HandleRequest(runtimeConfig)),
	}

	log.Printf("Started server on port: %d\n", runtimeConfig.Port)
	if runtimeConfig.UseHttps {
		log.Fatal(s.ListenAndServeTLS(runtimeConfig.CertFile, runtimeConfig.KeyFile))
	} else {
		log.Fatal(s.ListenAndServe())
	}
}
