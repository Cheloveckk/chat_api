package socket

import (
	"chat/api/internal/socket/errs"
	"chat/api/internal/task"
	"chat/api/pkgs/myjson"
	"context"
	"fmt"
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
		return errs.ErrClientChanClosed
	}
	select {
	case client.Send <- data:
		return nil
	default:
		return errs.ErrClientChanFull
	}
}

// Close user send channel
func (client *Client) CloseSend() {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.isClosed {
		return
	}
	close(client.Send)
	return
}

func (client *Client) CloseConn() {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.Conn.Close()
	return
}
func (client *Client) ReadWorker(ctx context.Context, tasks chan<- task.Task, unreg chan<- unregData) {
	var readerTask task.ReaderTask
	var err error
	defer log.Printf("Читающий поток пользователя connID: %d, userid: %d отключен\n", client.ConnID, client.UserID)
	for {
		err = client.Conn.ReadJSON(&readerTask)
		
		readerTask.ConnID = client.ConnID
		readerTask.UserID = client.UserID

		// Добавить обработку ошибок

		if err != nil {
			log.Println(err)
			if _, ok := err.(*websocket.CloseError); ok == true {
				log.Printf("Пользователь неактивен: %d", client.ConnID)
				select {
				case <-ctx.Done():
					client.CloseConn()
					return
				case unreg <- unregData{connID: readerTask.ConnID, userID: readerTask.UserID}:
				}
				return
			} else {
				client.HandleClientErr(unreg, fmt.Errorf("ReaderWorker: %w", myjson.ErrJSONFormat))
			}
		} else {
			select {
			case <-ctx.Done():
			case tasks <- task.Task{
				Type:   readerTask.Type,
				ConnID: readerTask.ConnID,
				UserID: readerTask.UserID,
				Data:   readerTask.Data,
			}:
			}
		}
	}

}

func (client *Client) WriteWorker(ctx context.Context) {
	defer log.Printf("Пишущий поток пользователя connID: %d, userid: %d отключен\n", client.ConnID, client.UserID)
	defer client.CloseConn()
	var err error
	for {
		select {
		case data := <-client.Send:
			if data == nil {
				return
			}
			err = client.Conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (client *Client) StartWork(h *Hub) {
	go client.ReadWorker(h.Ctx, h.Tasks, h.Unreg)
	go client.WriteWorker(h.Ctx)
}

func (client *Client) HandleClientErr(unreg chan<- unregData, err error) {
	err = client.SafeSend([]byte(err.Error()))
	if err != nil {
		unreg <- unregData{connID: client.ConnID, userID: client.UserID}
		client.CloseConn()
	}
	return
}
