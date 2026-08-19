package socket

import (
	"chat/api/internal/tasks"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ConnID int
	UserID int
	Conn   *websocket.Conn
	Send   chan []byte
}

func (client *Client) ReadWorker(ch chan tasks.Task) {
	var data tasks.Task
	var err error
	for {
		err = client.Conn.ReadJSON(&data)
		if _, ok := err.(*websocket.CloseError); ok == true {
			return
		} else if err != nil {
			log.Fatalln(err.Error())
		} else {
			ch <- data
		}
	}
}
func (client *Client) WriteWorker() {
	var err error
	for data := range client.Send {
		err = client.Conn.WriteJSON(data)
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
