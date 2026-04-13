package ws

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 512
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	UserID   uuid.UUID
	Email    string
	TokenExp time.Time
	send     chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uuid.UUID, email string, tokenExp time.Time) *Client {
	return &Client{
		Hub:      hub,
		Conn:     conn,
		UserID:   userID,
		Email:    email,
		TokenExp: tokenExp,
		send:     make(chan []byte, 64),
	}
}

func (c *Client) Send(data []byte) {
	select {
	case c.send <- data:
	default:
		// channel full, drop message
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMsgSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	tokenCheck := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		tokenCheck.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-tokenCheck.C:
			if time.Now().After(c.TokenExp) {
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				msg := []byte(`{"type":"force_logout","payload":{"reason":"token_expired"}}`)
				c.Conn.WriteMessage(websocket.TextMessage, msg)
				c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "token expired"))
				return
			}
		}
	}
}
