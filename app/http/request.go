package http

import (
	"bufio"
	"fmt"
	"net/textproto"
	"strings"

	"golang.org/x/net/http/httpguts"
)

// Request represents an HTTP request
type Request struct {
	Method string
	Path   string
	Proto  string
	Header Header
}

func badStringErr(what, val string) error { return fmt.Errorf("%s: %s", what, val) }

func ReadRequest(b *bufio.Reader) (req *Request, err error) {
	// textproto handle text which are basically in streams and parse accordingly with clrf
	tp := textproto.NewReader(b)
	req = new(Request)
	requestLine, err := tp.ReadLine()
	if err != nil {
		return nil, err
	}

	var ok bool
	req.Method, req.Path, req.Proto, ok = parseRequestLine(requestLine)
	if !ok {
		return nil, badStringErr("Malformed HTTP request", requestLine)
	}
	// validate method
	if valid := isValidMethod(req.Method); !valid {
		return nil, badStringErr("Malformed HTTP request", requestLine)
	}

	// PARSING HEADERs
	mineHeaders, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	req.Header = Header(mineHeaders)
	if len(req.Header["Host"]) > 1 {
		return nil, fmt.Errorf("too many Host in header")
	}

	// 	// 1. Create a byte slice of the exact size needed.
	// bodyBuffer := make([]byte, contentLength)

	// // 2. Use io.ReadFull to read from your bufio.Reader 'b'
	// //    and completely fill the bodyBuffer.
	// _, err := io.ReadFull(b, bodyBuffer)
	// if err != nil {
	//     // This can happen if the client closes the connection
	//     // or sends a body smaller than Content-Length.
	//     return nil, err
	// }

	// // 3. At this point, bodyBuffer holds the request body.
	// //    You can now assign it to your request struct.
	// req.Body = bodyBuffer

	return req, nil
}

// parse request line to method, uri, proto
func parseRequestLine(s string) (method, requestURI, proto string, ok bool) {
	method, rest, ok1 := strings.Cut(s, " ")
	requestURI, proto, ok2 := strings.Cut(rest, " ")
	if !ok1 || !ok2 {
		return "", "", "", false
	}
	return method, requestURI, proto, true
}

// according to HTTP spec, methods can be extended.
// the only restriction is that it should be valid token.
// for easier implementation, httpguts is used.
func isValidMethod(method string) bool {
	return httpguts.ValidHeaderFieldName(method)
}
