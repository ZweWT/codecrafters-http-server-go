package main

import (
	"log"
	"net"
	"os"
	"strings"
)

var (
	InfoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

var FileDirectory = "/temp/"

func getDirectoryFlag(args []string) (string, bool) {
	for i, arg := range args {
		if arg == "--directory" {
			return args[i+1], true
		}
	}
	return "", false
}

func main() {
	InfoLogger.Println("Logs from your program will appear here!")

	if filePath, ok := getDirectoryFlag(os.Args[1:]); ok {
		absolutePath, err := os.Getwd()
		if err != nil {
			ErrorLogger.Printf("error getting current directory: %s\n", err.Error())
		}
		InfoLogger.Printf("current path: %s\n", absolutePath)
		
		// If filePath is absolute, use it directly; otherwise, concatenate with current directory
		if strings.HasPrefix(filePath, "/") {
			FileDirectory = filePath
		} else {
			FileDirectory = absolutePath + "/" + filePath
		}
	}

	InfoLogger.Printf("directory: %s\n", FileDirectory)

	ln, err := net.Listen("tcp", "0.0.0.0:4221")
	InfoLogger.Println("serving from port 4221")
	if err != nil {
		ErrorLogger.Fatalln("Failed to bind to port 4221")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			ErrorLogger.Println("Failed to accept connection", err)
			continue
		}

		// other than net.conn, extra info such as file path where static file are stored
		// should be passed here. by this, dynamic file location can be provided to users.
		// am i over-thinking ?
		// yes, as a http server module, it should be doing one thing only.
		// (load-balancing can be introduced via others such as nginx)
		// it should be handling http request through a lifecycle
		go handleConnection(conn)
	}
}
