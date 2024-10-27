package proxy

import (
	"fmt"
	"io"
	"net/http"

	"github.com/patrickmn/go-cache"

	"github.com/sreeram-venkitesh/reverse-proxy/pkg/config"
)

func HandleRequest(cfg config.Config) http.HandlerFunc {
	proxyCache := cache.New(defaultExpiration, purgeTime)

	return func(rw http.ResponseWriter, r *http.Request) {
		// Based on incoming request we use the host name to find
		// the router and the service url it is pointing to from
		// the config.yaml file.
		targetRouterHost := r.Host

		var currentRouter config.Router

		// Go through the list of routers defined in config.yaml
		// and find the current requested router based on hostname
		for _, router := range cfg.Routers {
			if router.Host == targetRouterHost {
				currentRouter = router
			}
		}

		// Once we have the targeted router, we know which service
		// this router is pointing to. We can get the url of this service.
		serviceUrl, err := cfg.GetServiceUrl(currentRouter.Service)
		if err != nil {
			fmt.Printf("Error for router %s: %s", targetRouterHost, err)
		}

		// Once we get the service url, we will proxy our request to the url.
		targetUrl := fmt.Sprintf("%s%s", serviceUrl, r.URL.Path)

		cacheKey := generateCacheKey(r.Method, targetUrl, r.Header)

		// If response exists in cache, we read the data from cache and 
		// return the cached response to the client
		if cachedResponse, found := proxyCache.Get(cacheKey); found {
			response := cachedResponse.(*CachedResponse)

			for header, values := range response.Headers {
				for _, value := range values {
					rw.Header().Set(header, value)
				}
			}

			rw.Header().Set("X-Proxy-Cache", "HIT")
			rw.WriteHeader(response.StatusCode)

			rw.Write(response.Body)
			return
		}

		rw.Header().Set("X-Proxy-Cache", "MISS")

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

		body, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(rw, "Error reading response body", http.StatusInternalServerError)
			return
		}

		// Copying headers from target server's response to the proxy response
		for header, values := range res.Header {
			for _, value := range values {
				rw.Header().Set(header, value)
			}
		}

		if shouldCache(r.Method, res.StatusCode) {
			cachedResponse := &CachedResponse{
				Body:       body,
				Headers:    res.Header.Clone(),
				StatusCode: res.StatusCode,
			}
			proxyCache.Set(cacheKey, cachedResponse, defaultExpiration)
		}

		rw.WriteHeader(res.StatusCode)
		rw.Write(body)
	}
}
