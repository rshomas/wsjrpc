package wsjrpc

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	router       *Router
	upgrader     websocket.Upgrader
	onDisconnect func(addr string)
	clients      map[*websocket.Conn]string
	mu           sync.Mutex
}

func NewServer() *Server {
	return &Server{
		router:   NewRouter(),
		upgrader: websocket.Upgrader{},
		clients:  make(map[*websocket.Conn]string),
	}
}

func (s *Server) Use(mw Middleware) {
	s.router.Use(mw)
}

func (s *Server) Register(method string, handler HandlerFunc) {
	s.router.Register(method, handler)
}

func (s *Server) OnDisconnect(handler func(addr string)) {
	s.onDisconnect = handler
}

func (s *Server) Listen(addr string, endpoint string) error {
	http.HandleFunc(endpoint, s.handleWS)
	log.Printf("JSON-RPC WS server listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	addr := r.RemoteAddr
	s.mu.Lock()
	s.clients[conn] = addr
	s.mu.Unlock()
	log.Printf("Client connected: %s", addr)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Connection %s closed: %v", addr, err)
			if s.onDisconnect != nil {
				s.onDisconnect(addr)
			}
			s.mu.Lock()
			delete(s.clients, conn)
			s.mu.Unlock()
			break
		}

		var req Request
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("Invalid JSON:", err)
			s.sendError(conn, req.ID, -32700, "Invalid JSON")
			continue
		}

		if req.JSONRPC != "2.0" {
			log.Println("Invalid JSON-RPC version")
			s.sendError(conn, req.ID, -32600, "Invalid JSON-RPC version")
			continue
		}

		resp := s.router.Handle(&req)
		data, _ := json.Marshal(resp)
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (s *Server) sendError(conn *websocket.Conn, id any, code int, message string) {
	resp := Response{
		JSONRPC: "2.0",
		Error:   &Error{Code: code, Message: message},
		ID:      id,
	}
	conn.WriteJSON(resp)
}
