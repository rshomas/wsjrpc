package main

import (
	"encoding/json"
	"log"

	"github.com/rshomas/wsjrpc"
)

func main() {
	srv := wsjrpc.NewServer()

	// Middleware (logging)
	srv.Use(func(next wsjrpc.HandlerFunc) wsjrpc.HandlerFunc {
		return func(method string, params json.RawMessage) (any, *wsjrpc.Error) {
			log.Printf("Call %s with %s", method, string(params))
			return next(method, params)
		}
	})

	// Echo method
	srv.Register("echo", func(method string, params json.RawMessage) (any, *wsjrpc.Error) {
		var p struct{ Text string }
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, &wsjrpc.Error{Code: -32602, Message: "Invalid params"}
		}
		return p.Text, nil
	})

	// Disconnect event
	srv.OnDisconnect(func(addr string) {
		log.Printf("Client %s disconnected", addr)
	})

	log.Fatal(srv.Listen(":8080", "/ws"))
}
