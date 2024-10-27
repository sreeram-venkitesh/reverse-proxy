package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/sreeram-venkitesh/reverse-proxy/pkg/config"
)

func TestProxy(t *testing.T) {
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("Response: %s\n", string(body))

	if string(body) != "Hello world" {
		t.Errorf("Expected response body to include \"Hello World\" but got %s", string(body))
	} else {
		fmt.Println("Received the expected response from the proxy server")
	}
}