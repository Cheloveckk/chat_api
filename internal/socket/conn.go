package socket

import (
	"chat/api/internal/jwt"
	"chat/api/internal/socket/errs"
	"log"
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
		http.Error(w, "not authorizated", http.StatusUnauthorized)
		return nil, errs.ErrNotauthorized
	}
	ws, err := Upg.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	ws, err = SetSettings(ws)
	if err != nil {
		ws.Close()
		return nil, err
	}

	connID := GetNewID()
	log.Printf("Подключился user_id: %d, присвоена conn_id: %d", connID, userID)
	return NewClient(connID, userID, ws), nil
}
func SetSettings(ws *websocket.Conn) (*websocket.Conn, error) {
	return ws, nil
}
