package socket

import (
	"chat/api/internal/task"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Keep sharded  hubs
type HubMap struct {
	mu    sync.RWMutex
	Count int
	Hubs  []*Hub
}

// Keep UsersMap and let interact with clients in this
type Hub struct {
	mu       sync.RWMutex
	Tasks    chan task.Task
	Unreg    chan unregData
	Map      UsersMap
	isClosed bool
	wg       sync.WaitGroup
	Ctx      context.Context
	Cancel   context.CancelFunc
}

// Using UserID from database
type UsersMap map[int][]*Client

// Using connID
type ConnsMap map[int]*websocket.Conn

// Count - count of shards (hubs), cap - cap of shards
func NewHubMap(count int, cap int) *HubMap {
	slice := make([]*Hub, 0, count)
	for i := 0; i < count; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		h := &Hub{
			Tasks:  make(chan task.Task, 32),
			Map:    make(UsersMap, cap),
			Unreg:  make(chan unregData, 32),
			Ctx:    ctx,
			Cancel: cancel,
		}
		slice = append(slice, h)
		h.wg.Add(1)
		go unregWorker(h)
	}
	return &HubMap{
		Count: count,
		Hubs:  slice,
	}
}

// Create n workers for each hub
func (hubMap *HubMap) CreateWorkers(n int) {
	hubMap.mu.RLock()
	defer hubMap.mu.RUnlock()
	fmt.Println(hubMap.Hubs)
	for _, h := range hubMap.Hubs {
		for i := 0; i < n; i++ {
			h.wg.Add(1)
			go func(h *Hub) {
				for {
					defer h.wg.Done()
					select {
					case t, ok := <-h.Tasks:
						if ok != true {
							return
						}
						HandleRequest(hubMap, h, t)
					case <-h.Ctx.Done():
						return
					}
				}
			}(h)
		}
	}
}

// Get the Hub from HubMap using mutex
func (m *HubMap) GetHub(userID int) *Hub {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getHub(userID)
}

// Get the Hub from HubMap
func (m *HubMap) getHub(userID int) *Hub {
	shard := (userID) % m.Count
	h := m.Hubs[shard]
	return h
}

// Add the connection into the HubMap
func (m *HubMap) AddClient(client *Client) {
	h := m.GetHub(client.UserID)

	h.mu.Lock()
	defer h.mu.Unlock()
	if h.Map == nil {
		return
	}
	fmt.Println("до добавлени пользователя в срез в хабе", h.Map[client.UserID])
	h.Map[client.UserID] = append(h.Map[client.UserID], client)
	fmt.Println("после", h.Map[client.UserID])
}

func (m *HubMap) DeleteConn(connID int, userID int) {
	h := m.GetHub(userID)
	h.mu.Lock()

	h.deleteConn(connID, userID)

	h.mu.Unlock()
}

func (m *HubMap) GetConn(connID int, userID int) *Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h := m.getHub(userID)
	return h.GetConn(connID, userID)
}

func (h *Hub) GetConn(connID int, userID int) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.getConn(connID, userID)
}

func (h *Hub) getConn(connID int, userID int) *Client {
	if h.Map == nil {
		return nil
	}
	for _, cl := range h.Map[userID] {
		if cl.ConnID == connID {
			return cl
		}
	}
	return nil
}

func (h *Hub) DeleteConn(connID, userID int) {
	h.mu.Lock()
	h.deleteConn(connID, userID)
	h.mu.Unlock()
}

func (h *Hub) deleteConn(connID, userID int) {
	if h.Map == nil {
		return
	}
	slice := h.Map[userID]
	for pos, cl := range slice {
		log.Println("Сравнение id подключения для нахождения неактивного пользователя")
		if cl.ConnID == connID {
			cl.Close()
			cl.Conn.Close()

			l := len(h.Map[userID])
			slice[pos] = slice[l-1]
			slice[l-1] = nil
			h.Map[userID] = slice[:l-1]
			log.Printf("Удален пользователь с id = %d \n", connID)
		}
		if len(slice) == 0 {
			delete(h.Map, userID)
		}
		return
	}
}

// Don`t use it. Delete last hub of HubMap
func (hubMap *HubMap) delHub() {
	hubMap.mu.Lock()

	hubMap.Count--
	h := hubMap.Hubs[hubMap.Count]
	hubMap.mu.Unlock()
	h.mu.Lock()
	h.Cancel()
	h.mu.Unlock()
	h.wg.Add(1)
	go hubMap.evacuate(h)
	h.wg.Wait()
	close(h.Tasks)
	close(h.Unreg)
	hubMap.Hubs[hubMap.Count] = nil

}

// Don`t use it
func (hubMap *HubMap) evacuate(h *Hub) {
	if hubMap.Count == 0 {
		return
	}
	h.mu.Lock()
	var newHub *Hub
	for id, clients := range h.Map {
		newHub = hubMap.getHub(id)
		newHub.mu.Lock()
		newHub.Map[id] = clients
		newHub.mu.Unlock()
	}
	h.mu.Unlock()
	h.wg.Done()
	return
}
