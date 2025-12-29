package temperature

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
)

// MockHTTPClient is a mock implementation of HTTPDoer
type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestGetTemperatureByLocation(t *testing.T) {
	tests := []struct {
		name          string
		location      string
		apiKey        string
		mockResponse  string
		mockStatus    int
		mockErr       error
		expectError   bool
		expectedTempC float64
		expectedTempF float64
		expectedTempK float64
	}{
		{
			name:          "Success",
			location:      "London",
			apiKey:        "dummy-key",
			mockResponse:  `{"current": {"temp_c": 20.0, "temp_f": 68.0}}`,
			mockStatus:    http.StatusOK,
			mockErr:       nil,
			expectError:   false,
			expectedTempC: 20.0,
			expectedTempF: 68.0,
			expectedTempK: 293.0,
		},
		{
			name:         "Network Error",
			location:     "London",
			apiKey:       "dummy-key",
			mockResponse: "",
			mockStatus:   0,
			mockErr:      errors.New("network error"),
			expectError:  true,
		},
		{
			name:         "Invalid JSON",
			location:     "London",
			apiKey:       "dummy-key",
			mockResponse: `invalid-json`,
			mockStatus:   http.StatusOK,
			mockErr:      nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock response body
			body := io.NopCloser(bytes.NewReader([]byte(tt.mockResponse)))

			client := NewWeatherAPIClient(&MockHTTPClient{
				Response: &http.Response{
					StatusCode: tt.mockStatus,
					Body:       body,
				},
				Err: tt.mockErr,
			}, tt.apiKey)

			resp, err := client.GetTemperatureByLocation(context.Background(), tt.location)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if resp.TempC != tt.expectedTempC {
					t.Errorf("expected TempC %.2f, got %.2f", tt.expectedTempC, resp.TempC)
				}
				if resp.TempF != tt.expectedTempF {
					t.Errorf("expected TempF %.2f, got %.2f", tt.expectedTempF, resp.TempF)
				}
				if resp.TempK != tt.expectedTempK {
					t.Errorf("expected TempK %.2f, got %.2f", tt.expectedTempK, resp.TempK)
				}
			}
		})
	}
}
