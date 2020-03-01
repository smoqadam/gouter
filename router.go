package router

import (
	"context"
	"net/http"
)

type key int

const (
	contextKey key = iota
	varsKey
)

type Router struct {
	// Routes stores a collection of Route struct
	Routes []Route

	// ctx is an interface type will be accessible from http.request
	ctx    interface{}
}


// NewRouter return a new instance of Router
func NewRouter() *Router {
	return &Router{}
}

// GET register a GET request
func (r *Router) GET(path string, h http.HandlerFunc) *Route {
	return r.AddRoute(path, http.MethodGet, h)
}

// POST register a POST request
func (r *Router) POST(path string, h http.HandlerFunc) *Route {
	return r.AddRoute(path, http.MethodPost, h)
}

// PUT register a PUT request
func (r *Router) PUT(path string, h http.HandlerFunc) *Route {
	return r.AddRoute(path, http.MethodPut, h)
}

// PATCH register a PATCH request
func (r *Router) PATCH(path string, h http.HandlerFunc) *Route {
	return r.AddRoute(path, http.MethodPatch, h)
}

// DELETE register a DELETE request
func (r *Router) DELETE(path string, h http.HandlerFunc) *Route {
	return r.AddRoute(path, http.MethodDelete, h)
}

// AddRoute create a new Route and append it to Routes slice
func (r *Router) AddRoute(path string, method string, h http.HandlerFunc) *Route {
	route := NewRoute(path, method, h)
	r.Routes = append(r.Routes, route)
	return &route
}

// With send an interface along side the http.request.
// It is accessible with router.Context() function
func (r *Router) With(i interface{}) *Router {
	r.ctx = i
	return r
}

// ServeHTTP implement http.handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := context.WithValue(req.Context(), contextKey, r.ctx)
	req = req.WithContext(ctx)

	var match *Route
	var h http.Handler
	for _, route := range r.Routes {
		if route.Match(req) {
			vars := route.extractVars(req)
			ctx := context.WithValue(req.Context(), varsKey, vars)
			req = req.WithContext(ctx)
			match = &route
			break
		}
	}

	if match != nil && match.method != req.Method {
		h = &MethodNotAllowed{}
	}

	if h == nil && match != nil {
		h = match.dispatch()
	}

	if match == nil || h == nil {
		h = http.NotFoundHandler()
	}

	h.ServeHTTP(w, req)
}

// Vars return a map of variables defined on the route.
func Vars(req *http.Request) map[string]string {
	if v := req.Context().Value(varsKey); v != nil {
		return v.(map[string]string)
	}
	return nil
}

func Context(req *http.Request) interface{} {
	if v := req.Context().Value(contextKey); v != nil {
		return v
	}
	return nil
}
