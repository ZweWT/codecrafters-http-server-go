package main

import (
	"net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	InfoLogger.Println("Logs from your program will appear here!")

	ln, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		ErrorLogger.Fatalln("Failed to bind to port 4221")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			ErrorLogger.Println("Failed to accept connection", err)
			continue
		}

		go handleConnection(conn)
	}
}
