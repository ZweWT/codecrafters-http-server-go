package main

import "net"

// handleConnection processes a single client connection
func handleConnection(conn net.Conn) error {
	// handle connection
	defer conn.Close()
	InfoLogger.Printf("New Connection from %s", conn.RemoteAddr())
	defer InfoLogger.Printf("Connection closed from %s", conn.RemoteAddr())
	var writeErr error

	request, err := requestParser(conn)
	if err != nil {
		ErrorLogger.Printf("failed to parsed the request")
		return err
	}

	InfoLogger.Printf("request_received method=%s path=%s remote=%s",
		request.Method, request.Path, conn.RemoteAddr())

	if request.Method == "GET" && request.Path == "/" {
		writeErr = writeResponse(conn, 200, "OK")
	} else {
		writeErr = writeResponse(conn, 404, "Not Found")
	}

	if writeErr != nil {
		ErrorLogger.Printf("Connection closed with error: %v", writeErr)
		return writeErr
	}

	return nil
}
