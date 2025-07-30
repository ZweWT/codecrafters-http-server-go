package main

import (
	"bufio"
	"strings"
	"testing"
)

// request parser
//
//	possible input are the one provided the code crafter excercise.
func TestRequestParser(t *testing.T) {
	input := "GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	req, err := requestParser(reader)

	if err != nil {
		t.Fatalf("Expect no error, got %v", err)
	}

	if req.Method != "GET" {
		t.Fatalf("Expected method GET, got %s", req.Method)
	}

	if req.Path != "/index.html" {
		t.Fatalf("Expected path /index.html , got %s", req.Path)
	}

	if req.Proto != "HTTP/1.1" {
		t.Fatalf("Expected proto HTTP/1.1, got %s", req.Proto)
	}

	invalidinput := "GET\r\n\r\n"
	invalidReader := bufio.NewReader(strings.NewReader(invalidinput))

	_, err = requestParser(invalidReader)
	if err == nil {
		t.Errorf("Expected error for invalid request, got nil")
	}
}
