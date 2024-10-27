package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name           string
		configFilePath string
		configData     string
		expectedError  string
	}{
		{
			name: "valid config.yaml file",
			configData: `
port: 8080

routers:
  - host: iamfoo.localhost:8080
    service: whoami
  - host: testing.localhost:8080
    service: whoami
  - host: sreeram.localhost:8080
    service: whoami

services:
  - name: whoami
    url: "http://localhost:9000"
`,
		},
		{
			name: "invalid config.yaml file",
			configData: `
invalid file
`,
			expectedError: "Error parsing config.yaml",
		},
		{
			name:           "invalid config.yaml file name",
			configFilePath: "wrongconfig.yaml",
			configData: `
port: 8080
`,
			expectedError: "Error reading config.yaml",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			tmpDir := t.TempDir()

			configPath := filepath.Join(tmpDir, "config.yaml")

			err := os.WriteFile(configPath, []byte(test.configData), 0644)
			if err != nil {
				t.Fatalf("Failed to create test config file: %v", err)
			}

			if test.configFilePath != "" {
				configPath = test.configFilePath
			}

			_, err = loadConfig(configPath)

			if test.expectedError != "" {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), test.expectedError) {
					fmt.Println("Error is")
					fmt.Println(err.Error())
					fmt.Println("We want to have")
					fmt.Println(test.expectedError)
					t.Errorf("Expected error to contain %v, but got %v", test.expectedError, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
