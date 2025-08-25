package main

import (
	"encoding/json"
	"log"

	"github.com/rshomas/wsjrpc"
)

func main() {
	srv := wsjrpc.NewServer()

	// echo
	srv.Register("echo", func(method string, params json.RawMessage) (any, *wsjrpc.Error) {
		var p struct{ Text string }
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, &wsjrpc.Error{Code: -32602, Message: "Invalid params"}
		}
		return p.Text, nil
	})

	log.Fatal(srv.Listen(":8080"))
}
