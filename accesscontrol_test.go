package accesscontrol_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccessControl_Whitelist(t *testing.T) {
	config := &Config{
		Whitelist: []string{"192.168.1.1"},
		HeaderKey: "X-Test-Header",
		HeaderValue: "TestValue",
	}

	nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	handler, err := New(context.Background(), nextHandler, config, "test-plugin")
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		expectCode int
	}{
		{
			name:       "Allowed IP",
			remoteAddr: "192.168.1.1:12345",
			headers:    nil,
			expectCode: http.StatusOK,
		},
		{
			name:       "Blocked IP",
			remoteAddr: "10.0.0.1:12345",
			headers:    nil,
			expectCode: http.StatusForbidden,
		},
		{
			name:       "Allowed Header",
			remoteAddr: "10.0.0.1:12345",
			headers:    map[string]string{"X-Test-Header": "TestValue"},
			expectCode: http.StatusOK,
		},
		{
			name:       "Blocked Header",
			remoteAddr: "10.0.0.1:12345",
			headers:    map[string]string{"X-Test-Header": "InvalidValue"},
			expectCode: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			req.RemoteAddr = test.remoteAddr
			for key, value := range test.headers {
				req.Header.Set(key, value)
			}

			rw := httptest.NewRecorder()
			handler.ServeHTTP(rw, req)

			if rw.Code != test.expectCode {
				t.Errorf("Expected status code %d, got %d", test.expectCode, rw.Code)
			}
		})
	}
}
