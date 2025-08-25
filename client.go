package wsjrpc

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	closed atomic.Bool
	reqID  atomic.Uint64
}

func NewClient(url string) (*Client, error) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &Client{conn: c}, nil
}

func (c *Client) Call(method string, params any, result any) error {
	id := c.reqID.Add(1)

	// Формируем запрос
	req := Request{
		JSONRPC: "2.0",
		Method:  method,
		ID:      id,
	}
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return err
		}
		req.Params = b
	}

	data, _ := json.Marshal(req)
	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}

	// Ждём ответ
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return err
	}

	var resp Response
	if err := json.Unmarshal(msg, &resp); err != nil {
		return err
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	if result != nil {
		b, _ := json.Marshal(resp.Result)
		return json.Unmarshal(b, result)
	}
	return nil
}

func (c *Client) Close() error {
	if c.closed.Swap(true) {
		return nil
	}
	return c.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
}

func (c *Client) Ping() error {
	deadline := time.Now().Add(5 * time.Second)
	return c.conn.WriteControl(websocket.PingMessage, []byte("ping"), deadline)
}
