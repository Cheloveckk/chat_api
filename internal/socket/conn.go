package socket

import (
	jwt "chat/api/internal/JWT"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upg = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func InitConn(w http.ResponseWriter, r *http.Request) (*Client, error) {
	userID := jwt.GetID(r)
	if userID == 0 {
		return nil, errors.New("not autharisated")
	}
	ws, err := Upg.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	ws, err = SetSettings(ws)
	if err != nil {
		return nil, err
	}
	return NewClient(GetNewID(), int(userID), ws), err
}
func SetSettings(ws *websocket.Conn) (*websocket.Conn, error) {
	return ws, nil
}
