# wsjrpc

A lightweight Go module for **JSON-RPC 2.0 over WebSocket**.
It provides request routing, middleware, disconnect events, and a simple synchronous client.

## Features

- JSON-RPC 2.0 client/server over WebSocket
- Synchronous WebSocket client (`Call`)
- Method routing via `Register`
- Middleware support for logging, authentication, etc.
- `OnDisconnect` callback for connection lifecycle handling

## Usage

### Client

```go
  client, err := wsjrpc.NewClient("ws://localhost:8080/ws")
  if err != nil {
    log.Fatal("connect:", err)
  }
  defer client.Close()

  var resp string
  if err := client.Call("echo", map[string]string{"text": "Hello from client"}, &resp); err != nil {
    log.Fatal("rpc error:", err)
  }
  fmt.Println("Server response:", resp)
```

### Server

```go
  srv := wsjrpc.NewServer()

  // Middleware (logging)
  srv.Use(func(next wsjrpc.HandlerFunc) wsjrpc.HandlerFunc {
    return func(method string, params json.RawMessage) (any, *wsjrpc.Error) {
      log.Printf("Call %s with %s", method, string(params))
      return next(method, params)
    }
  })

  // 'echo' method
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

  log.Fatal(srv.Listen(":8080", "/ws))
```

## Roadmap

- Support for notifications (no id)

- Asynchronous calls (CallAsync)

- Broadcast to all clients

- Batch requests

## License

`MIT`
