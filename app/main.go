package main

import (
	"log"
	"net"
	"os"
)

var (
	InfoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
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
