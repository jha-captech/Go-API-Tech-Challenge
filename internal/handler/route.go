package handler

import "net/http"

type Route interface {
	http.Handler
	Pattern() string
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
