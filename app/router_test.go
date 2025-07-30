package main

import (
	"strings"
	"testing"
)

// MockResponseWriter implements ResponseWriter for testing
type MockResponseWriter struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       string
	WriteError error
}

// NewMockResponseWriter creates a new mock response writer
func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    make(map[string]string),
	}
}

// SetStatus sets the status code and text
func (m *MockResponseWriter) SetStatus(code int, text string) {
	m.StatusCode = code
	m.StatusText = text
}

// SetHeader sets a header
func (m *MockResponseWriter) SetHeader(key, value string) {
	m.Headers[key] = value
}

// Write records the response but doesn't actually write anywhere
func (m *MockResponseWriter) Write() error {
	return m.WriteError
}

// GetBody returns the response body
func (m *MockResponseWriter) GetBody() string {
	return m.Body
}

func (m *MockResponseWriter) SetBody(body string) {
	m.Body = body
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	router.HandlePrefix("/echo/", func(rw ResponseWriter, req *Request) error {
		echoText := strings.TrimPrefix(req.Path, "/echo/")
		rw.SetStatus(200, "OK")
		rw.SetBody(echoText)
		return rw.Write()
	})

	tests := []struct {
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{"/echo/abc", 200, "abc"},
		{"/echo/world", 200, "world"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := &Request{
				Method: "GET",
				Path:   tt.path,
				Proto:  "HTTP/1.1",
			}

			rw := NewMockResponseWriter()

			err := router.ServeHTTP(req, rw)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if rw.GetBody() != tt.expectedBody {
				t.Errorf("Expected boyd %q, got %q", tt.expectedBody, rw.GetBody())
			}

		})
	}

}
