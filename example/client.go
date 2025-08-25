package main

import (
	"fmt"
	"log"

	"github.com/rshomas/wsjrpc"
)

func main() {
	client, err := wsjrpc.NewClient("ws://localhost:8080/ws")
	if err != nil {
		log.Fatal("connect:", err)
	}
	defer client.Close()

	var resp string
	if err := client.Call("echo", map[string]string{"text": "Hello from Go client"}, &resp); err != nil {
		log.Fatal("rpc call:", err)
	}
	fmt.Println("Server response:", resp)
}
