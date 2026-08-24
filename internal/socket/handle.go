package socket

import (
	"chat/api/internal/task"
	"log"
)

var TaskMap map[string]TaskFunc = map[string]TaskFunc{
	"sendto": HandleSendTo,
}

type TaskFunc func(*HubMap, *Hub, task.Task) error

func HandleRequest(hubMap *HubMap, h *Hub, t task.Task) {
	f := TaskMap[t.Type]
	if f == nil {
		HandleError(h, t, ErrInvalidCommand)
		return
	}
	err := f(hubMap, h, t)
	if err != nil {
		return
	}
	HandleError(h, t, err)
}

func HandleError(h *Hub, t task.Task, err error) {
	if err == nil {
		return
	}
	cl := h.GetConn(t.ConnID, t.UserID)
	if cl == nil {
		return
	}
	err = cl.SafeSend([]byte(err.Error()))
	if err != nil {
		log.Printf("Пользователь неактивен: %d", t.ConnID)
		select {
		case h.Unreg <- unregData{t.ConnID, t.UserID}:
		case <-h.Ctx.Done():
		default:
			cl.Close()
			cl.Conn.Close()
		}
	}
	return
}
