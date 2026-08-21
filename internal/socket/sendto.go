package socket

import (
	"chat/api/internal/task"
	"encoding/json"
	"errors"
	"fmt"
)

type SendTo struct {
	To      int `json:"to"`
	UserID  int
	ConnID  int64
	Message string `json:"message"`
}

func ParseSendTo(data task.Task) (SendTo, error) {
	var s SendTo
	err := json.Unmarshal(data.Data, &s)
	if err != nil || s.To == 0 {
		return s, errors.New("incorrect data for command")
	}
	s.UserID = data.UserID
	s.ConnID = data.ConnID
	fmt.Println("отправить к ", s.To, "отправить от", s.UserID)
	return s, err
}

func HandleSendTo(hubMap *HubMap, currentHub *Hub, t task.Task) error {
	s, err := ParseSendTo(t)
	if err != nil {
		return err
	}

	currentHub.mu.RLock()
	defer currentHub.mu.RUnlock()

	h := hubMap.GetHub(s.To)
	if h == nil || h.Map == nil {
		return errors.New("Server connection error")
	}
	slice := h.Map[s.To]
	fmt.Println(slice)
	for _, cl := range slice {
		cl.Send <- []byte(s.Message)
	}
	return nil
}
