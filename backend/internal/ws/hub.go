package ws

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type Event struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload,omitempty"`
	ProjectID string      `json:"project_id,omitempty"`
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
	logger  *slog.Logger
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		logger:  logger,
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
	h.logger.Info("ws client connected", "user_id", c.UserID)
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
	h.logger.Info("ws client disconnected", "user_id", c.UserID)
}

func (h *Hub) BroadcastToUser(userID uuid.UUID, evt Event) {
	data, err := json.Marshal(evt)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.UserID == userID {
			c.Send(data)
		}
	}
}

func (h *Hub) BroadcastToAll(evt Event) {
	data, err := json.Marshal(evt)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		c.Send(data)
	}
}
