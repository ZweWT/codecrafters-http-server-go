package http

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"

	"golang.org/x/net/http/httpguts"
)

const MAX_BODY_SIZE = 1024 * 1024

// Request represents an HTTP request
type Request struct {
	Method string
	Path   string
	Proto  string
	Header Header
	Body   []byte
}

func badStringErr(what, val string) error { return fmt.Errorf("%s: %s", what, val) }

var ErrBodyTooLarge = fmt.Errorf("http: request body too large")

type maxByteReader struct {
	r io.Reader // underlying reader(bufio)
	n int64     // bytes remaining allowed
}

// instead of count up approach(checking till max_body_size), this func uses preventive truncation.
// this logic allow any attempt to read to further for violation ErrBodyTooLarge
func (l *maxByteReader) Read(p []byte) (n int, err error) {
	// violation check
	if l.n <= 0 {
		return 0, io.EOF
	}

	// if the caller asks more than n left, truncate till n(preventive truncation)
	// why this is great?
	// imagine l.n is 10 bytes, this prevents underlying reader from reading 11th byte.
	// also ensures the limit is hit exactly at boundary.
	if int64(len(p)) > l.n {
		p = p[:l.n]
	}

	n, err = l.r.Read(p)
	l.n -= int64(n)
	return
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

	contentLength := req.Header.Get("Content-Length")
	contentLengthInt, _ := strconv.Atoi(contentLength)
	fmt.Printf("content length: %v and max body size: %v\n", contentLengthInt, MAX_BODY_SIZE)
	if contentLengthInt > MAX_BODY_SIZE {
		return nil, ErrBodyTooLarge
	}

	if contentLengthInt > 0 {
		limitedReader := &maxByteReader{
			r: b,
			n: int64(contentLengthInt),
		}

		buffer, err := io.ReadAll(limitedReader)
		if err != nil {
			fmt.Printf("error reading with limited reader: %s", err.Error())
			return nil, err
		}
		req.Body = buffer
	}

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
