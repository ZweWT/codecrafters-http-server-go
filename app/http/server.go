package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sort"
	"strings"
)

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}

type ServeMux struct {
	m  map[string]muxEntry
	es []muxEntry // sorted from longest to shortest for prefix routes
}

type muxEntry struct {
	h       Handler
	pattern string
}

func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	h, _ := mux.findHandler(r)
	fmt.Printf("found handler: %v\n", h)
	if h == nil {
		w.SetStatus(404, "Not Found")
		w.SetBody([]byte("Not Found"))
		w.Write()
		return
	}
	h.ServeHTTP(w, r)
}

func (mux *ServeMux) findHandler(r *Request) (h Handler, pattern string) {
	path := r.Path
	// exact keyword match
	fmt.Printf("before keyword match for path finding: %s\n", path)
	v, ok := mux.m[path]
	fmt.Printf("found in keyword match: %t\n", ok)
	if ok {
		return v.h, v.pattern
	}

	for _, e := range mux.es {
		fmt.Printf("matching with register route: %s\n", e.pattern)
		// matches the longest parts first
		if strings.HasPrefix(path, e.pattern) {
			return e.h, e.pattern
		}
	}

	return nil, ""
}

func (mux *ServeMux) Handle(pattern string, handler Handler) {
	if _, exist := mux.m[pattern]; exist {
		panic("multiple registration for same routes")
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}

	e := muxEntry{
		h:       handler,
		pattern: pattern,
	}
	mux.m[pattern] = e //single keyword matches

	// matches with prefix
	// prefix the routes ends in /, i.e /echo/
	if len(pattern) > 1 && pattern[len(pattern)-1] == '/' {
		mux.es = appendSorted(mux.es, e)
	}

}
func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
	n := len(es)

	i := sort.Search(n, func(i int) bool {
		return len(es[i].pattern) <= len(e.pattern)
	})

	if i == n {
		return append(es, e)
	}

	// we already know i points to where we should insert
	// so first, grow the size of slice
	// move the shorter entries down
	// and insert into the i index
	es = append(es, muxEntry{})
	copy(es[i+1:], es[i:])
	es[i] = e
	return es
}

func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("nil handler")
	}

	mux.Handle(pattern, HandlerFunc(handler))
}

func NewServeMux() *ServeMux {
	return &ServeMux{}
}

var defaultServeMux ServeMux
var DefaultServeMux = &defaultServeMux

// this allow to set defaultServeMux
type serverHandler struct {
	svr *Server
}

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.svr.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	handler.ServeHTTP(rw, req)
}

type Server struct {
	Addr    string
	Handler Handler
}

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to bind port %s", addr)
		return err
	}
	return s.Serve(ln)
}

func (s *Server) Serve(ln net.Listener) error {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				continue
			}
			return err
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	defer conn.Close()

	b := bufio.NewReader(conn)

	for {
		req, err := ReadRequest(b)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Printf("error reading request: %s", err.Error())
			res := NewResponse(conn, req)
			if err == ErrBodyTooLarge {
				res.SetStatus(413, "Payload Too Large")
				res.SetBody([]byte("Payload Too Large"))
			} else {
				res.SetStatus(400, "Bad Request")
				res.SetBody([]byte("Bad Request"))
			}
			return res.Write()
		}

		res := NewResponse(conn, req)
		serverHandler{svr: s}.ServeHTTP(res, req)

		if strings.ToLower(req.Header.Get("Connection")) == "close" {
			return nil
		}
	}
}

func ListenAndServe(addr string, handler Handler) error {
	s := &Server{
		Addr:    addr,
		Handler: handler,
	}

	return s.ListenAndServe()
}

func Serve(ln net.Listener, handler Handler) error {
	s := &Server{
		Handler: handler,
	}
	return s.Serve(ln)
}
