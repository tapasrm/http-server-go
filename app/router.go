
package main

import (
	"net/http"
	"strings"
)

type Handler func(req *Request) *Response

type Route struct {
	Path    string
	Method  string
	Handler Handler
}

type Router struct {
	routes []*Route
}

func NewRouter() *Router {
	return &Router{
		routes: make([]*Route, 0),
	}
}

func (r *Router) Add(method, path string, handler Handler) {
	route := &Route{
		Path:    path,
		Method:  method,
		Handler: handler,
	}
	r.routes = append(r.routes, route)
}

func (r *Router) Find(req *Request) *Response {
	for _, route := range r.routes {
		if route.Method != req.Method {
			continue
		}

		routeParts := strings.Split(route.Path, "/")
		requestParts := strings.Split(req.Path, "/")

		if len(routeParts) != len(requestParts) {
			continue
		}

		params := make(map[string]string)
		match := true
		for i, part := range routeParts {
			if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
				paramName := part[1 : len(part)-1]
				params[paramName] = requestParts[i]
			} else if part != requestParts[i] {
				match = false
				break
			}
		}

		if match {
			req.Params = params
			return route.Handler(req)
		}
	}

	return &Response{
		StatusCode: http.StatusNotFound,
		Body:       "Not Found",
	}
}
