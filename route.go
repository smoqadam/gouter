package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Middleware func(handler http.Handler) http.Handler

type Route struct {
	path    string
	name    string
	handler http.Handler
	method  string
	mw      []Middleware
	where   map[string]string
	vars    map[string]string
}

// NewRoute create a new route
func NewRoute(path string, method string, handler http.HandlerFunc) Route {
	return Route{
		path:    path,
		handler: handler,
		method:  method,
		vars:    make(map[string]string),
		where:   make(map[string]string),
	}
}

// Name assign a name for the route
func (r *Route) Name(s string) *Route {
	r.name = s
	return r
}

// Match return true if the requested path would match with the current route path
func (r *Route) Match(req *http.Request) bool {
	regex := regexp.MustCompile(`{([^}]*)}`)
	matches := regex.FindAllStringSubmatch(r.path, -1)
	p := r.path
	for _, v := range matches {
		s := fmt.Sprintf("{%s}", v[1])
		p = strings.Replace(p, s, r.where[v[1]], -1)
	}
	regex, err := regexp.Compile(p)
	if err != nil {
		return false
	}
	matches = regex.FindAllStringSubmatch(req.URL.Path, -1)
	for _, match := range matches {

		if regex.Match([]byte(match[0])) {
			return true
		}
	}
	return false
}

func (r *Route) clear(s string) string {
	s = strings.Replace(s, "{", "", -1)
	s = strings.Replace(s, "}", "", -1)
	return s
}

// Where define a regex pattern for the variables in the route path
func (r *Route) Where(key string, pattern string) *Route {
	r.where[key] = fmt.Sprintf("(%s)", pattern)
	return r
}

// Middleware register a collection of middleware functions and sort them
func (r *Route) Middleware(mw ...Middleware) *Route {
	r.mw = mw

	//TODO: Fix this
	for i := len(r.mw)/2 - 1; i >= 0; i-- {
		opp := len(r.mw) - 1 - i
		r.mw[i], r.mw[opp] = r.mw[opp], r.mw[i]
	}
	return r
}

// extractVars parse the requested URL and return key/value pair of
// variables defined in the route path.
func (r *Route) extractVars(req *http.Request) map[string]string {
	url := strings.Split(req.URL.Path, "/")
	path := strings.Split(r.clear(r.path), "/")
	vars := make(map[string]string)
	for i := 0; i < len(url); i++ {
		if _, ok := r.where[path[i]]; ok {
			vars[path[i]] = url[i]
		}
	}
	return vars
}

// dispatch run route middleswares if any then run the route handler
func (r *Route) dispatch() http.Handler {
	for _, m := range r.mw {
		r.handler = m(r.handler)
	}
	return r.handler
}
