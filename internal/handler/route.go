package handler

import "net/http"

type Route interface {
	Handler() http.HandlerFunc
	Pattern() string
}

func NewRoute(pattern string, handler http.HandlerFunc) Route {
	return routeImpl{
		handler: handler,
		pattern: pattern,
	}
}

type routeImpl struct {
	handler http.HandlerFunc
	pattern string
}

func (s routeImpl) Pattern() string {
	return s.pattern
}

func (s routeImpl) Handler() http.HandlerFunc {
	return s.handler
}

// type RouteImpl struct {
// 	handler http.Handler
// 	pattern string
// }

// func (s RouteImpl) Pattern() string {
// 	return s.pattern
// }

// func (s RouteImpl) Handler() http.Handler {
// 	return s.handler
// }

// func newRoute(pattern string, handler http.Handler) Route {
// 	return RouteImpl{
// 		handler: handler,
// 		pattern: pattern,
// 	}
// }
