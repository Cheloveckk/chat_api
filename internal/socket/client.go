package socket

import (
	"chat/api/internal/task"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ConnID int64
	UserID int
	Conn   *websocket.Conn
	Send   chan []byte
}

func NewClient(connID int64, userID int, ws *websocket.Conn) *Client {
	return &Client{
		ConnID: connID,
		UserID: userID,
		Conn:   ws,
		Send:   make(chan []byte, 16),
	}
}
func (client *Client) ReadWorker(ch chan task.Task) {
	var data task.Task
	var err error
	for {
		err = client.Conn.ReadJSON(&data)
		data.ConnID = client.ConnID
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok == true {
				return
			} else {
				client.Send <- []byte("json request error")
			}
		} else {
			ch <- data
		}
	}
}
func (client *Client) WriteWorker() {
	var err error
	for data := range client.Send {
		if data == nil {
			return
		}
		err = client.Conn.WriteMessage(websocket.BinaryMessage, data)
		if _, ok := err.(*websocket.CloseError); ok == true {
			return
		} else if err != nil {
			log.Fatalln(err.Error())
		}
	}
}
func (client *Client) StartWork(ch chan task.Task) {
	go client.ReadWorker(ch)
	go client.WriteWorker()
}
