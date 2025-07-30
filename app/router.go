package main

import "strings"

// Handler is a function that handles an HTTP request
type HandlerFunc func(res ResponseWriter, resquest *Request) error

type Router struct {
	routes       map[string]HandlerFunc
	prefixRoutes []PrefixRoute
}

type PrefixRoute struct {
	prefix  string
	handler HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes:       make(map[string]HandlerFunc),
		prefixRoutes: []PrefixRoute{},
	}
}

// Handle registers the handler for route without prefix.
func (r *Router) Handle(path string, h HandlerFunc) {
	r.routes[path] = h
}

// HandlePrefix registers the handler for route with prefix such as /echo/abc
func (r *Router) HandlePrefix(prefix string, h HandlerFunc) {
	r.prefixRoutes = append(r.prefixRoutes, PrefixRoute{
		prefix:  prefix,
		handler: h,
	})
}

func (r *Router) FindHandler(req *Request) (HandlerFunc, bool) {
	// exact keyword match
	// does not work for prefix such as /echo/abc
	if h, ok := r.routes[req.Path]; ok {
		return h, true
	}

	// in order echo/abc to work we need to fall-back and execute.
	for _, route := range r.prefixRoutes {
		if strings.HasPrefix(req.Path, route.prefix) {
			return route.handler, true
		}
	}
	return nil, false
}

// processes a http request and write a ResponseWriter
func (r *Router) ServeHTTP(req *Request, rw ResponseWriter) error {
	h, found := r.FindHandler(req)
	if !found {
		rw.SetStatus(404, "Not Found")
		return rw.Write()
	}

	return h(rw, req)
}
