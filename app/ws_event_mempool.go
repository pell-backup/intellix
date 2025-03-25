package app

import (
	"context"
	"cosmossdk.io/log"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
)

var _ mempool.ExtMempool = (*WsEventMempool)(nil)

// WsEventMempool implements the mempool.ExtMempool interface and embeds a WebSocket hub
type WsEventMempool struct {
	hub    *MempoolEventHub
	logger log.Logger
}

func NewWsEventMempool(hub *MempoolEventHub, logger log.Logger) *WsEventMempool {
	return &WsEventMempool{
		hub:    hub,
		logger: logger,
	}
}

// Insert is called when a new transaction is inserted into the mempool.
// It logs the transaction details and broadcasts an event to all connected WebSocket clients.
func (m *WsEventMempool) Insert(ctx context.Context, tx sdk.Tx) error {
	m.logger.Info("New tx inserted into mempool")

	// Log the messages contained in the transaction
	msgs, err := tx.GetMsgsV2()
	if err != nil {
		return err
	}

	var jsonMsgs []json.RawMessage
	for _, msg := range msgs {
		b, err := protojson.Marshal(msg)
		if err != nil {
			return err
		}
		jsonMsgs = append(jsonMsgs, b)
	}

	// Marshal the messages to JSON
	jsonBytes, err := json.Marshal(jsonMsgs)
	if err != nil {
		return err
	}

	// Broadcast the event message to all connected WebSocket clients
	m.hub.broadcast <- jsonBytes

	// Return an error or nil based on business logic
	return nil
}

// Other interface methods can be implemented as needed

func (WsEventMempool) Select(context.Context, [][]byte) mempool.Iterator     { return nil }
func (WsEventMempool) SelectBy(context.Context, [][]byte, func(sdk.Tx) bool) {}
func (WsEventMempool) CountTx() int                                          { return 0 }
func (WsEventMempool) Remove(sdk.Tx) error                                   { return nil }

// -------------------- WebSocket hub & MempoolWSClient implementation --------------------

// MempoolEventHub manages all WebSocket client connections and handles broadcast messages.
type MempoolEventHub struct {
	clients    map[*MempoolWSClient]bool
	broadcast  chan []byte
	register   chan *MempoolWSClient
	unregister chan *MempoolWSClient
}

// NewMempoolEventHub initializes and returns a new MempoolEventHub instance.
func NewMempoolEventHub() *MempoolEventHub {
	return &MempoolEventHub{
		clients:    make(map[*MempoolWSClient]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *MempoolWSClient),
		unregister: make(chan *MempoolWSClient),
	}
}

// run continuously handles client registration, unregistration, and broadcasting messages.
func (h *MempoolEventHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// MempoolWSClient represents a WebSocket client connection.
type MempoolWSClient struct {
	hub  *MempoolEventHub
	conn *websocket.Conn
	send chan []byte
}

// readPump handles incoming messages from the WebSocket connection.
// In this example, it simply reads messages and unregisters the client on error.
func (c *MempoolWSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
		// If needed, process the message from the client here
	}
}

// writePump sends messages from the hub to the WebSocket client.
func (c *MempoolWSClient) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// If the channel is closed, send a close message and exit
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	// Allow all origins; in production, validate the request origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// serveMempoolEventWs upgrades the HTTP connection to a WebSocket connection and creates a new MempoolWSClient.
func serveMempoolEventWs(hub *MempoolEventHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	client := &MempoolWSClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	hub.register <- client

	// Start goroutines for writing and reading messages
	go client.writePump()
	go client.readPump()
}
