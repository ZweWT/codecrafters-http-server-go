package main

import (
	"bufio"
	"fmt"
	"strings"
)

// requestParser reads and parses an HTTP request from a connection
func requestParser(reader *bufio.Reader) (*Request, error) {
	// Read the first line of the HTTP request
	// TODO: remove this request line and replace with readline from below.
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
	r := &Request{
		Method:  parsedRequestLine[0],
		Path:    parsedRequestLine[1],
		Proto:   parsedRequestLine[2],
		Headers: make(map[string]string),
	}

	for {
		line, err := readLine(reader)
		if err != nil {
			ErrorLogger.Printf("Error reading headers: %v", err)
			return nil, err
		}
		// Empty line signals the end of headers
		if line == "" {
			break
		}
		InfoLogger.Printf("header after line break %s", line)

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// add into request.Header for each key and value
			r.Headers[key] = value
		}
		InfoLogger.Printf("Header: %v", r.Headers)
	}

	return r, nil
}

// readLine reads a line from a bufio.Reader
func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	InfoLogger.Printf("readLine func only parse for headers: %s", line)

	// Trim the trailing CRLF
	return strings.TrimSuffix(line, "\r\n"), nil
}

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
