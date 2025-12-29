package buscacep

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

func TestGetLocationByCep(t *testing.T) {
	tests := []struct {
		name           string
		cep            string
		mockResponse   string
		mockStatusCode int
		mockErr        error
		expectedLoc    string
		expectErr      bool
	}{
		{
			name:           "Success",
			cep:            "89035400",
			mockResponse:   `{"localidade": "Blumenau"}`,
			mockStatusCode: http.StatusOK,
			mockErr:        nil,
			expectedLoc:    "Blumenau",
			expectErr:      false,
		},
		{
			name:           "Network Error",
			cep:            "89035400",
			mockResponse:   "",
			mockStatusCode: 0,
			mockErr:        errors.New("network error"),
			expectedLoc:    "",
			expectErr:      true,
		},
		{
			name:           "Invalid JSON",
			cep:            "89035400",
			mockResponse:   `{invalid json`,
			mockStatusCode: http.StatusOK,
			mockErr:        nil,
			expectedLoc:    "",
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock response body
			body := io.NopCloser(bytes.NewReader([]byte(tt.mockResponse)))

			// Setup the mock client
			mockClient := &MockHTTPClient{
				Response: &http.Response{
					StatusCode: tt.mockStatusCode,
					Body:       body,
				},
				Err: tt.mockErr,
			}

			// Initialize the service with the mock client
			service := NewViaCepClient(mockClient)

			// Call the method under test
			loc, err := service.GetLocationByCep(context.Background(), tt.cep)

			// Assertions
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if loc != tt.expectedLoc {
					t.Errorf("expected location %s, got %s", tt.expectedLoc, loc)
				}
			}
		})
	}
}
