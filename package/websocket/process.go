package websocket

import (
	"errors"
	"sync"
)

type (
	Request struct {
		Client *Client
		Send   []byte
	}
	DisposeFunc    func(*Request) *Response
	MiddlewareFunc func(DisposeFunc) DisposeFunc

	Router struct {
		handlers sync.Map
	}

	Route struct {
		final       DisposeFunc
		middlewares []MiddlewareFunc
	}
)

var (
	RouteManager *Router
)

type DistributeHandler interface {
	Distribute(client *Client, message []byte) (err error)
}

func NewRouter() *Router {
	RouteManager = &Router{}
	return RouteManager
}

// Register event to memory
func (r *Router) Register(event string, handler DisposeFunc, middlewareFunc ...MiddlewareFunc) {
	route := &Route{
		final: handler,
	}
	for _, middleware := range middlewareFunc {
		route.Use(middleware)
	}
	r.handlers.Store(event, route)
}

// GetRoute Get register's route
func (r *Router) GetRoute(event string) (route *Route, err error) {
	value, ok := r.handlers.Load(event)
	if !ok {
		return nil, errors.New("current event not supported")
	}
	route, okk := value.(*Route)
	if !okk {
		return nil, errors.New("handler type error")
	}
	return
}

// Use Route Create a new middleware chain
func (r *Route) Use(middleware MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}

// Execute the middleware chain
func (r *Route) Execute(request *Request) *Response {
	handler := r.final
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	return handler(request)
}
