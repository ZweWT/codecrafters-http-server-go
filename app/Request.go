package main

// Request represents an HTTP request
type Request struct {
	Method  string
	Path    string
	Proto   string
	Headers map[string]string
}
