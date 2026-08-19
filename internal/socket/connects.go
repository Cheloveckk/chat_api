package socket

import (
	"chat/api/internal/tasks"
	"sync"

	"github.com/gorilla/websocket"
)

// Keep sharded  hubs
type HubMap struct {
	Count int
	Hubs  []*Hub
}

// Keep UsersMap and let interact with clients in this
type Hub struct {
	mu    sync.RWMutex
	tasks chan tasks.Task
	Map   UsersMap
}

// Using connID
type ConnsMap map[int]*websocket.Conn

// Using UserID from database
type UsersMap map[int][]*Client

// Count - count of shards (hubs), len - len of shards
func NewHubMap(count int, len int) *HubMap {
	slice := make([]*Hub, 0, count)
	for i := 0; i < count; i++ {
		slice = append(slice, &Hub{
			tasks: make(chan tasks.Task, 32),
			Map:   make(UsersMap, len),
		})
	}
	return &HubMap{
		Count: count,
		Hubs:  slice,
	}
}

// Get the Hub from HubMap
func (m *HubMap) GetHub(userID int) *Hub {
	shard := userID % m.Count
	h := m.Hubs[shard]
	return h
}

// Add the connection into the HubMap
func (m *HubMap) AddConn(connID int, userID int, ws *websocket.Conn) {
	h := m.GetHub(userID)

	h.mu.Lock()
	defer h.mu.Unlock()
	h.Map[userID] = append(h.Map[userID], &Client{
		ConnID: connID,
		UserID: userID,
		Conn:   ws,
		Send:   make(chan []byte, 16),
	})
}

func (m *HubMap) DeleteConn(connID int, userID int) {
	h := m.GetHub(userID)

	h.mu.Lock()
	defer h.mu.Unlock()
	for pos, cl := range h.Map[userID] {
		if cl.ConnID == connID {
			close(cl.Send)
			l := len(h.Map[userID])
			h.Map[userID][pos] = h.Map[userID][l-1]
			h.Map[userID][l-1] = nil
		}
	}
}

func (m *HubMap) GetClient(connID int, userID int) *Client {
	h := m.GetHub(userID)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, cl := range h.Map[userID] {
		if cl.ConnID == connID {
			return cl
		}
	}
	return nil
}
