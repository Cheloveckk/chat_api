package socket

import (
	"chat/api/internal/task"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Keep sharded  hubs
type HubMap struct {
	Count int
	Hubs  []*Hub
}

// Create n workers for each hub
func (hubMap *HubMap) CreateWorkers(n int) {
	fmt.Println(hubMap.Hubs)
	for _, h := range hubMap.Hubs {
		for i := 0; i < n; i++ {
			go func(h *Hub) {
				for t := range h.Ch {
					err := HandleRequest(hubMap, h, t)
					if err != nil {
						HandleError(h, t, err)
					}
				}
			}(h)
		}

	}
}

// Keep UsersMap and let interact with clients in this
type Hub struct {
	mu  sync.RWMutex
	Ch  chan task.Task
	Map UsersMap
}

// Using UserID from database
type UsersMap map[int][]*Client

// Using connID
type ConnsMap map[int64]*websocket.Conn

// Count - count of shards (hubs), len - len of shards
func NewHubMap(count int, len int) *HubMap {
	slice := make([]*Hub, 0, count)
	for i := 0; i < count; i++ {
		slice = append(slice, &Hub{
			Ch:  make(chan task.Task, 32),
			Map: make(UsersMap, len),
		})
	}
	return &HubMap{
		Count: count,
		Hubs:  slice,
	}
}

// Get the Hub from HubMap
func (m *HubMap) GetHub(userID int) *Hub {
	fmt.Println(userID, userID-1)
	shard := (userID) % m.Count
	fmt.Println(shard)

	h := m.Hubs[shard]
	return h
}

// Add the connection into the HubMap
func (m *HubMap) AddClient(client *Client) {
	h := m.GetHub(client.UserID)

	h.mu.Lock()
	defer h.mu.Unlock()
	fmt.Println("до добавлени пользователя в срез в хабе", h.Map[client.UserID])
	h.Map[client.UserID] = append(h.Map[client.UserID], client)
	fmt.Println("после", h.Map[client.UserID])
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
	return h.GetClient(connID, userID)

}

func (h *Hub) GetClient(connID int64, userID int) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, cl := range h.Map[userID] {
		if cl.ConnID == connID {
			return cl
		}
	}
	return nil
}
