package main

import (
	"bufio"
	"net"
	"strings"
)

// handleConnection processes a single client connection
func handleConnection(conn net.Conn) error {
	// handle connection
	defer conn.Close()
	InfoLogger.Printf("New Connection from %s", conn.RemoteAddr())
	defer InfoLogger.Printf("Connection closed from %s", conn.RemoteAddr())

	// Create a new buffered reader from the connection
	r := bufio.NewReader(conn)

	// Parse the HTTP request
	request, err := requestParser(r)
	if err != nil {
		ErrorLogger.Printf("Failed to parse the request: %v", err)
		// Send a 400 Bad Request response
		return writeResponse(conn, 400, "Bad Request")
	}

	// Log the received request
	InfoLogger.Printf("Request received: method=%s path=%s protocol=%s remote=%s",
		request.Method, request.Path, request.Proto, conn.RemoteAddr())

	// Skip the request headers
	for {
		line, err := readLine(r)
		if err != nil {
			ErrorLogger.Printf("Error reading headers: %v", err)
			return err
		}
		// Empty line signals the end of headers
		if line == "" {
			break
		}
		InfoLogger.Printf("Header: %s", line)
	}

	// Handle the request based on the method and path
	var writeErr error
	if request.Method == "GET" {
		if request.Path == "/" {
			// Root path returns a simple response
			writeErr = writeResponse(conn, 200, "")
		} else if strings.HasPrefix(request.Path, "/echo/") {
			// Echo path returns the path component after /echo/
			echoText := strings.TrimPrefix(request.Path, "/echo/")
			writeErr = writeResponse(conn, 200, echoText)
		} else {
			// Any other path returns a 404 Not Found
			writeErr = writeResponse(conn, 404, "Not Found")
		}
	} else {
		// Non-GET methods return a 405 Method Not Allowed
		writeErr = writeResponse(conn, 405, "Method Not Allowed")
	}

	if writeErr != nil {
		ErrorLogger.Printf("Error writing response: %v", writeErr)
		return writeErr
	}

	return nil
}
