package socket

import (
	"chat/api/internal/task"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ConnID   int
	UserID   int
	Conn     *websocket.Conn
	mu       sync.Mutex
	isClosed bool
	Send     chan []byte
}

func NewClient(connID int, userID int, ws *websocket.Conn) *Client {
	return &Client{
		ConnID: connID,
		UserID: userID,
		Conn:   ws,
		Send:   make(chan []byte, 16),
	}
}

func (client *Client) SafeSend(data []byte) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.isClosed {
		return ErrClientChanClosed
	}
	select {
	case client.Send <- data:
		log.Println(string(data))
		return nil
	default:
		return ErrClientChanFull
	}
}

func (client *Client) Close() {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.isClosed {
		return
	}
	close(client.Send)
	return
}
func (client *Client) ReadWorker(h *Hub) {
	defer h.wg.Done()
	var task task.Task
	var err error
	for {
		err = client.Conn.ReadJSON(&task)
		task.ConnID = client.ConnID
		task.UserID = client.UserID
		h.mu.RLock()
		defer h.mu.RUnlock()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok == true {
				log.Printf("Пользователь неактивен: %d", client.ConnID)
				select {
				case <-h.Ctx.Done():
					return
				case h.Unreg <- unregData{connID: task.ConnID, userID: task.UserID}:

				}
				return
			} else {
				client.HandleClientErr(h, ErrJSONFormat)
			}
		} else {
			select {
			case <-h.Ctx.Done():
				return
			case h.Tasks <- task:
			}
		}
	}

}

func (client *Client) WriteWorker(h *Hub) {
	defer h.wg.Done()
	var err error
	for data := range client.Send {
		if data == nil {
			return
		}
		err = client.Conn.WriteJSON(string(data))
		if err != nil {
			client.Conn.Close()
			return
		}
	}
}

func (client *Client) StartWork(h *Hub) {
	go client.ReadWorker(h)
	go client.WriteWorker(h)
}

func (client *Client) HandleClientErr(h *Hub, err error) {
	err = client.SafeSend([]byte(err.Error()))
	if err != nil {
		if err == ErrClientChanClosed {
			return
		}
		select {
		case <-h.Ctx.Done():
		case h.Unreg <- unregData{client.ConnID, client.UserID}:
		}
	}
	return
}
