package config

import (
	"testing"
)

func TestGetServiceUrl(t *testing.T) {
	testCases := []struct {
		name          string
		config        Config
		serviceName   string
		expectedUrl   string
		expectedError string
	}{
		{
			name: "valid service exists in config",
			config: Config{
				Services: []Service{
					{Name: "iamfoo", URL: "http://localhost:9000"},
				},
			},
			serviceName: "iamfoo",
			expectedUrl: "http://localhost:9000",
		},
		{
			name: "service doesn't exist in config",
			config: Config{
				Services: []Service{
					{Name: "iamfoo", URL: "http://localhost:9000"},
				},
			},
			serviceName:   "iamfoo2",
			expectedUrl:   "",
			expectedError: "Service not found\n",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			url, err := test.config.GetServiceUrl(test.serviceName)

			if test.expectedError != "" {
				if err == nil {
					t.Error("Expected error but got none")
				} else if err.Error() != test.expectedError {
					t.Errorf("Expected error %v, but got %v", test.expectedError, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if url != test.expectedUrl {
				t.Errorf("Expected %v, got %v", test.expectedUrl, url)
			}
		})
	}
}
