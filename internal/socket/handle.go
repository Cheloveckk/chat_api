package socket

import (
	"chat/api/internal/task"
	"errors"
)

var TaskMap map[string]TaskFunc = map[string]TaskFunc{
	"sendto": HandleSendTo,
}

type TaskFunc func(*HubMap, *Hub, task.Task) error

func HandleRequest(hubMap *HubMap, currentHub *Hub, t task.Task) error {
	f := TaskMap[t.Type]
	if f == nil {
		return errors.New("command doesn`t exist")
	}
	return f(hubMap, currentHub, t)
}

func HandleError(currentHub *Hub, t task.Task, err error) {
	currentHub.mu.RLock()
	cl := currentHub.GetClient(t.ConnID, t.UserID)
	select {
	case cl.Send <- []byte(err.Error()):
		return
	default:
		cl.Conn.Close()
	}
}
