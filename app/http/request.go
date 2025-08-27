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
		return nil, fmt.Errorf("%s %q", "Malformed HTTP request", requestLine)
	}
	// validate method
	if valid := isValidMethod(req.Method); !valid {
		// return 400 bad request error
	}
	// parse proto and validate the proto major and minor
	//
	// PARSING HEADERs
	mineHeaders, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	req.Header = Header(mineHeaders)

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
