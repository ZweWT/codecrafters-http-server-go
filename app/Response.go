package main

import (
	"fmt"
	"net"
)

type ResponseWriter interface {
	SetStatus(code int, text string)
	SetHeader(key, value string)
	SetBody(body string)
	GetBody() string
	Write() error
}

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       string
	conn       net.Conn
}

func NewResponse(conn net.Conn) *Response {
	return &Response{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    make(map[string]string),
		conn:       conn,
	}
}

// SetStatus sets the status code and text
func (r *Response) SetStatus(code int, text string) {
	r.StatusCode = code
	r.StatusText = text
}

// SetHeader sets a header in the response
func (r *Response) SetHeader(key, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
}

// GetBody returns the response body
func (r *Response) GetBody() string {
	return r.Body
}

func (r *Response) SetBody(body string) {
	r.Body = body
}

func (r *Response) Write() error {

	if _, ok := r.Headers["Content-Type"]; !ok {
		r.Headers["Content-Type"] = "text/plain"
	}

	r.Headers["Content-Length"] = fmt.Sprintf("%d", len(r.Body))

	// Build response string
	responseString := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, r.StatusText)

	// Add headers
	for key, value := range r.Headers {
		responseString += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	// Add empty line and body
	responseString += "\r\n" + r.Body

	// Write to connection
	_, err := r.conn.Write([]byte(responseString))
	return err
}
