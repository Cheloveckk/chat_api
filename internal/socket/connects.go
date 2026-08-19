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
	mu      sync.RWMutex
	ChTasks chan tasks.Task
	Map     UsersMap
}

// Using connID
type ConnsMap map[int64]*websocket.Conn

// Using UserID from database
type UsersMap map[int][]*Client

// Count - count of shards (hubs), len - len of shards
func NewHubMap(count int, len int) *HubMap {
	slice := make([]*Hub, 0, count)
	for i := 0; i < count; i++ {
		slice = append(slice, &Hub{
			ChTasks: make(chan tasks.Task, 32),
			Map:     make(UsersMap, len),
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
func (m *HubMap) AddClient(client *Client) {
	h := m.GetHub(client.UserID)

	h.mu.Lock()
	defer h.mu.Unlock()

	h.Map[client.UserID] = append(h.Map[client.UserID], client)
}

func (m *HubMap) DeleteConn(connID int64, userID int) {
	h := m.GetHub(userID)

	h.mu.Lock()
	defer h.mu.Unlock()

	for pos, cl := range h.Map[userID] {
		if cl.ConnID == connID {
			close(cl.Send)
			l := len(h.Map[userID])
			h.Map[userID][pos] = h.Map[userID][l-1]
			h.Map[userID][l-1] = nil
		} else {
			delete(h.Map, userID)
		}
	}
}

func (m *HubMap) GetClient(connID int64, userID int) *Client {
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
