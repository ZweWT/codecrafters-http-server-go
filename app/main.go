package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
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

	serveMux := registerServeMux()
	server := http.Server{
		Addr:    ":4221",
		Handler: serveMux,
	}

	fmt.Printf("server mux : %v", serveMux)

	log.Fatal(server.ListenAndServe())
}

func registerServeMux() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.SetStatus(200, "OK")
		w.SetBody([]byte(""))
		w.Write()
	})

	serveMux.HandleFunc("/echo/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("echo route: %s", r.Path)
		echoText := strings.TrimPrefix(r.Path, "/echo/")
		w.SetStatus(200, "OK")
		w.SetBody([]byte(echoText))
		w.Write()
	})

	serveMux.HandleFunc("/echo/david", func(w http.ResponseWriter, r *http.Request) {
		w.SetStatus(200, "OK")
		w.SetBody([]byte("this is registerd echo david route"))
		w.Write()
	})

	serveMux.HandleFunc("/user-agent", func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		w.SetStatus(200, "OK")
		w.SetBody([]byte(userAgent))
		w.Write()
	})

	serveMux.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		fileName := strings.TrimPrefix(r.Path, "/files/")
		path := fmt.Sprintf("%s%s", FileDirectory, fileName)
		fmt.Printf("path: %s", path)
		if r.Method == "POST" {

			// os.WriteFile(path, []byte("hello"), os.ModeDevice.Perm())
			// w.SetStatus(200, "OK")
			// w.SetHeader("Content-Length", strconv.Itoa(len(contents)))
			// w.SetHeader("Content-Type", "application/octet-stream")
			// w.SetBody([]byte(contents))
			// w.Write()
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			ErrorLogger.Printf("reading file error: %s", err.Error())
			w.SetStatus(404, "Not Found")
			w.Write()
		}

		w.SetStatus(200, "OK")
		w.SetHeader("Content-Length", strconv.Itoa(len(contents)))
		w.SetHeader("Content-Type", "application/octet-stream")
		w.SetBody([]byte(contents))
		w.Write()
	})

	return serveMux
}
