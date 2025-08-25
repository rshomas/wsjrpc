package wsjrpc

import "encoding/json"

type HandlerFunc func(method string, params json.RawMessage) (any, *Error)
type Middleware func(HandlerFunc) HandlerFunc

type Router struct {
	handlers    map[string]HandlerFunc
	middlewares []Middleware
}

func NewRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc)}
}

func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) Register(method string, handler HandlerFunc) {
	h := handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}
	r.handlers[method] = func(_ string, params json.RawMessage) (any, *Error) {
		return h(method, params)
	}
}

func (r *Router) Handle(req *Request) *Response {
	h, ok := r.handlers[req.Method]
	if !ok {
		return &Response{JSONRPC: "2.0", Error: &Error{-32601, "Method not found"}, ID: req.ID}
	}
	result, err := h(req.Method, req.Params)
	if err != nil {
		return &Response{JSONRPC: "2.0", Error: err, ID: req.ID}
	}
	return &Response{JSONRPC: "2.0", Result: result, ID: req.ID}
}
