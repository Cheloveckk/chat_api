package socket

import (
	"chat/api/internal/tasks"
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
func (client *Client) ReadWorker(ch chan tasks.Task) {
	var data tasks.Task
	var err error
	for {
		err = client.Conn.ReadJSON(&data)
		data.From = client.ConnID
		if _, ok := err.(*websocket.CloseError); ok == true {
			return
		} else if err != nil {
			client.Send <- []byte("request error")
		} else {
			ch <- data
		}
	}
}
func (client *Client) WriteWorker() {
	var err error
	for data := range client.Send {
		err = client.Conn.WriteMessage(websocket.BinaryMessage, data)
		if _, ok := err.(*websocket.CloseError); ok == true {
			return
		} else if err != nil {
			log.Fatalln(err.Error())
		}
	}
}
func (client *Client) StartWork(ch chan tasks.Task) {
	go client.ReadWorker(ch)
	go client.WriteWorker()
}
