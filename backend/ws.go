package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSMessage map[string]interface{}

// ServeWS upgrades HTTP to WebSocket and creates a client
func ServeWS(gm *GameManager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	client := NewClient(conn, gm)
	client.Listen()
}

// Client represents a player connection
type Client struct {
	conn     *websocket.Conn
	gm       *GameManager
	username string
	gameId   string
	send     chan WSMessage
}

func NewClient(conn *websocket.Conn, gm *GameManager) *Client {
	return &Client{
		conn: conn,
		gm:   gm,
		send: make(chan WSMessage, 16),
	}
}

func (c *Client) Listen() {
	defer func() {
		c.conn.Close()
		if c.username != "" {
			c.gm.HandleDisconnect(c.username, c.gameId)
		}
	}()

	// Writer goroutine
	go func() {
		for msg := range c.send {
			c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := c.conn.WriteJSON(msg); err != nil {
				log.Println("write json err:", err)
				return
			}
		}
	}()

	// Reader loop
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read err:", err)
			return
		}
		var m WSMessage
		if err := json.Unmarshal(raw, &m); err != nil {
			log.Println("invalid json:", err)
			c.send <- WSMessage{"type": "error", "message": "invalid json"}
			continue
		}
		c.handleMessage(m)
	}
}

// handleMessage processes incoming messages from client
func (c *Client) handleMessage(m WSMessage) {
	t, _ := m["type"].(string)
	switch t {
	case "join":
		username, _ := m["username"].(string)
		c.username = username
		c.gm.JoinQueue(username, c)

	case "drop":
		colf, _ := m["column"].(float64)
		col := int(colf)
		gameId, _ := m["gameId"].(string)
		c.gm.HandleMove(gameId, c.username, col)

	case "rematch_request":
		gameId, _ := m["gameId"].(string)
		c.gm.HandleRematchRequest(gameId, c.username)

	default:
		c.send <- WSMessage{"type": "error", "message": "unknown type"}
	}
}
