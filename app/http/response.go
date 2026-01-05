package http

import (
	"fmt"
	"net"
	"strings"
)

type ResponseWriter interface {
	SetStatus(code int, text string)
	SetHeader(key, value string)
	SetBody(body []byte)
	GetBody() []byte
	Write() error
}

// there is no reason for user to use Response type, as responseWriter will be used.
type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
	conn       net.Conn
}

func NewResponse(conn net.Conn, req *Request) *Response {
	res := &Response{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    make(map[string]string),
		conn:       conn,
	}

	if req != nil && strings.ToLower(req.Header.Get("Connection")) == "close" {
		res.SetHeader("Connection", "close")
	} else {
		res.SetHeader("Connection", "keep-alive")
	}

	return res
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
func (r *Response) GetBody() []byte {
	return r.Body
}

func (r *Response) SetBody(body []byte) {
	r.Body = body
}

func (r *Response) Write() error {

	if _, ok := r.Headers["Content-Type"]; !ok {
		r.SetHeader("Content-Type", "text/plain")
	}

	if _, ok := r.Headers["Connection"]; !ok {
		r.SetHeader("Connection", "keep-alive")
	}

	r.SetHeader("Content-Length", fmt.Sprintf("%d", len(r.Body)))

	// Build response string
	headerString := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, r.StatusText)

	// Add headers
	for key, value := range r.Headers {
		headerString += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	// Add empty line and body
	headerString += "\r\n"

	responseString := append([]byte(headerString), r.Body...)

	// Write to connection
	_, err := r.conn.Write([]byte(responseString))
	if err != nil {
		return err
	}
	return nil
}
