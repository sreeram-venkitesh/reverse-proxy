package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/sreeram-venkitesh/reverse-proxy/pkg/config"
	"github.com/sreeram-venkitesh/reverse-proxy/pkg/proxy"
)

func loadConfig(configPath string) (config.Config, error) {
	var runtimeConfig config.Config

	// Read the config file
	yamlConfig, err := os.ReadFile(configPath)
	if err != nil {
		return runtimeConfig, fmt.Errorf("Error reading config.yaml: %v", err)
	}

	err = yaml.Unmarshal(yamlConfig, &runtimeConfig)
	if err != nil {
		return runtimeConfig, fmt.Errorf("Error parsing config.yaml: %v", err)
	}

	return runtimeConfig, nil
}

func main() {
	log.Printf("Reverse proxy is starting...")

	runtimeConfig, err := loadConfig("config.yaml")
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
	log.Fatal(s.ListenAndServe())
}