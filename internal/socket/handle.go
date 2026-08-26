package socket

import (
	"chat/api/internal/socket/errs"
	"chat/api/internal/task"
	"errors"
	"fmt"
)

var TaskMap map[string]TaskFunc = map[string]TaskFunc{
	"sendto":   HandleSendTo,
	"sendchat": HandleSendChat,
}

type TaskFunc func(*HubMap, *Hub, task.Task) error

func HandleRequest(hubMap *HubMap, h *Hub, t task.Task) {
	switch t.Type {
	case "delivery":
		HandleDeliverySendChat(h, t.Data)

	default:
		f := TaskMap[t.Type]
		if f == nil {
			HandleError(h, t, errs.ErrInvalidCommand)
			return
		}
		err := f(hubMap, h, t)
		if err == nil {
			return
		}
		HandleError(h, t, fmt.Errorf("Handle request: %w", err))
		return
	}
}

func HandleError(h *Hub, t task.Task, err error) {
	if errors.Is(err, nil) {
		return
	}
	cl := h.GetConn(t.ConnID, t.UserID)
	if cl == nil {
		return
	}
	cl.HandleClientErr(h.Unreg, err)
	return
}
