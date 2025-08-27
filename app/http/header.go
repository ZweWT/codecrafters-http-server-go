package http

import "net/textproto"

// Header represents the key-value pair in an HTTP header
type Header map[string][]string

func (h Header) Get(key string) string {
	if v := h[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}
