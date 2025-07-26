package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Request represents an HTTP request
type Request struct {
	Method string
	Path   string
	Proto  string
	// Headers map[string]string
}

// requestParser reads and parses an HTTP request from a connection
func requestParser(reader *bufio.Reader) (*Request, error) {
	// Read the first line of the HTTP request
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading request line: %w", err)
	}

	// Trim the trailing CRLF
	requestLine = strings.TrimSuffix(requestLine, "\r\n")
	InfoLogger.Printf("Request line: %s", requestLine)

	// Parse the request line into its components
	parsedRequestLine := strings.Fields(requestLine)
	if len(parsedRequestLine) < 3 {
		return nil, fmt.Errorf("invalid request line: %s", requestLine)
	}

	// Create and return a new Request object
	return &Request{
		Method: parsedRequestLine[0],
		Path:   parsedRequestLine[1],
		Proto:  parsedRequestLine[2],
	}, nil
}

// readLine reads a line from a bufio.Reader
func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Trim the trailing CRLF
	return strings.TrimSuffix(line, "\r\n"), nil
}

// writeResponse writes an HTTP response to a connection
func writeResponse(conn net.Conn, statusCode int, body string) error {
	statusText := statusTextForCode(statusCode)
	responseString := fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		statusCode, statusText, len(body), body)
	_, err := conn.Write([]byte(responseString))
	return err
}

// statusTextForCode returns the status text for a given HTTP status code
func statusTextForCode(code int) string {
	switch code {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}
