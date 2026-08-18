package socket

import "github.com/gorilla/websocket"

type User struct {
	ConnID int
	UserID int
	Conn   *websocket.Conn
}
