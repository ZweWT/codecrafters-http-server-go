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
func requestParser(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	// first line after '\n' occurance, which means it is only request line, no header
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	//parse requestLine
	parsedRequestLine := strings.Fields(requestLine)

	//parse string to request type
	return &Request{
		Method: parsedRequestLine[0],
		Path:   parsedRequestLine[1],
		Proto:  parsedRequestLine[2],
	}, nil
}

// writeResponse writes an HTTP response to a connection
func writeResponse(conn net.Conn, statusCode int, body string) error {
	responseString := fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", statusCode, body)
	_, err := conn.Write([]byte(responseString))
	return err
}
