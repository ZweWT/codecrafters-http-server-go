package http

import "testing"

var parseRequestLineTest = []struct {
	line, method, path, proto string
}{
	{"GET /index.html HTTP/1.1", "GET", "/index.html", "HTTP/1.1"},
	{"POST /submit?foo=bar HTTP/1.1", "POST", "/submit?foo=bar", "HTTP/1.1"},
	{"HEAD / HTTP/1.1", "HEAD", "/", "HTTP/1.1"},
}

func TestParseRequestLine(t *testing.T) {
	for i, tt := range parseRequestLineTest {
		method, uri, proto, ok := parseRequestLine(tt.line)
		if !ok {
			t.Errorf("#%d: FAILED TO SPLIT WITH SINGLE SPACE", i)
		}
		if method != tt.method {
			t.Errorf("#%d: gotMethod: %q wantMethod: %q ", i, method, tt.method)
		}
		if uri != tt.path {
			t.Errorf("#%d: gotURI: %q wantURI: %q ", i, uri, tt.path)
		}
		if proto != tt.proto {
			t.Errorf("#%d: gotProto: %q wantProto: %q ", i, proto, tt.proto)
		}
	}
}

var parseRequestErrorTest = []struct{}{}

func TestParseRequestError(t *testing.T) {
	for i, tt := range parseRequestErrorTest {
	}
}
