package main

import (
	"bufio"
	"net"
	"strings"
)

// router singleton
var routerInstance *Router

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
	res := NewResponse(conn)
	if err != nil {
		ErrorLogger.Printf("Failed to parse the request: %v", err)
		// Send a 400 Bad Request response
		res.SetStatus(400, "Bad Request")
		res.SetBody("Bad Request")
		return res.Write()
	}

	// Log the received request
	InfoLogger.Printf("Request received: method=%s path=%s protocol=%s remote=%s",
		request.Method, request.Path, request.Proto, conn.RemoteAddr())

	router := getRouter()
	err = router.ServeHTTP(request, res)

	if err != nil {
		ErrorLogger.Printf("Error serving HTTP: %v", err)
		return err
	}

	return nil
}

func getRouter() *Router {
	if routerInstance == nil {
		routerInstance = NewRouter()

		routerInstance.Handle("/", func(res ResponseWriter, req *Request) error {
			res.SetStatus(200, "OK")
			res.SetBody("")
			return res.Write()
		})

		// Register echo handler
		routerInstance.HandlePrefix("/echo/", func(res ResponseWriter, req *Request) error {
			echoText := strings.TrimPrefix(req.Path, "/echo/")
			res.SetStatus(200, "OK")
			res.SetBody(echoText)
			return res.Write()
		})

		routerInstance.Handle("/user-agent", func(res ResponseWriter, req *Request) error {
			userAgent := req.Headers["User-Agent"]
			res.SetStatus(200, "OK")
			res.SetBody(userAgent)
			return res.Write()
		})
	}

	return routerInstance
}
