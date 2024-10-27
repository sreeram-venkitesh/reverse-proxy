package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/sreeram-venkitesh/reverse-proxy/pkg/config"
)

func TestCache(t *testing.T) {
	dummyServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello world"))
	}))

	defer dummyServer.Close()

	dummyServerUrl, err := url.Parse(dummyServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	configData := fmt.Sprintf(`
port: 8080

routers:
  - host: dummy.localhost:8080
    service: dummy

services:
  - name: dummy
    url: "%s"
`, dummyServerUrl.String())

	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yaml")

	err = os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	runtimeConfig, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}

	testProxyHandler := HandleRequest(runtimeConfig)
	testProxyServer := httptest.NewUnstartedServer(testProxyHandler)
	listener, err := net.Listen("tcp", "localhost:8080")
	testProxyServer.Listener = listener
	testProxyServer.Start()
	defer testProxyServer.Close()

	testProxyServerClient := testProxyServer.Client()

	request, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	request.Host = "dummy.localhost:8080"

	res, err := testProxyServerClient.Do(request)
	defer res.Body.Close()

	if res.Header.Get("X-Proxy-Cache") != "MISS" {
		t.Errorf("Expected X-Proxy-Cache to be MISS, but got %s", res.Header.Get("X-Proxy-Cache"))
	} else {
		fmt.Println("X-Proxy-Cache is MISS in the first response as expected")
	}

	res, err = testProxyServerClient.Do(request)
	defer res.Body.Close()

	if res.Header.Get("X-Proxy-Cache") != "HIT" {
		t.Errorf("Expected X-Proxy-Cache to be HIT, but got %s", res.Header.Get("X-Proxy-Cache"))
	} else {
		fmt.Println("X-Proxy-Cache is HIT in the second response as expected")
	}
}
