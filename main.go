package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Definitions for the config file
// which defines routers and services
type Config struct {
	Port     int       `yaml:"port"`
	Routers  []Router  `yaml:"routers"`
	Services []Service `yaml:"services"`
}

type Service struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Router struct {
	Host    string `yaml:"host"`
	Service string `yaml:"service"`
}

var config Config

func (c *Config) GetServiceUrl(name string) (string, error) {
	for _, svc := range c.Services {
		if svc.Name == name {
			return svc.URL, nil
		}
	}
	return "", fmt.Errorf("Service not found\n")
}

func handleRequest(rw http.ResponseWriter, r *http.Request) {
	// Based on incoming request we use the host name to find 
	// the router and the service url it is pointing to from
	// the config.yaml file.
	targetRouterHost := r.Host

	var currentRouter Router

	// Go through the list of routers defined in config.yaml 
	// and find the current requested router based on hostname
	for _, router := range config.Routers {
		if router.Host == targetRouterHost {
			currentRouter = router
		}
	}

	// Once we have the targeted router, we know which service
	// this router is pointing to. We can get the url of this service.
	serviceUrl, err := config.GetServiceUrl(currentRouter.Service)
	if err != nil {
		fmt.Printf("Error for router %s: %s", targetRouterHost, err)
	}

	// Once we get the service url, we will proxy our request to the url.
	targetUrl := fmt.Sprintf("%s%s", serviceUrl, r.URL.Path)
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

	// Copying headers from target server's response to the proxy response
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

	// Read the config file
	yamlConfig, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config.yaml: %v", err)
	}

	err = yaml.Unmarshal(yamlConfig, &config)
	if err != nil {
		log.Fatalf("Error parsing config.yaml: %v", err)
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        http.HandlerFunc(handleRequest),
	}

	log.Printf("Started server on port: %d\n", config.Port)
	log.Fatal(s.ListenAndServe())
}
