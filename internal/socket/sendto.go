package socket

import (
	"chat/api/internal/socket/errs"
	"chat/api/internal/task"
	"chat/api/pkgs/myjson"
	"encoding/json"
	"fmt"
	"log"
)

type SendTo struct {
	To      int `json:"to"`
	UserID  int
	ConnID  int
	Message string `json:"message"`
}

func ParseSendTo(task task.Task) (SendTo, error) {
	d, ok := task.Data.(json.RawMessage)
	if !ok {
		return SendTo{}, fmt.Errorf("ParseSendTo: %w", myjson.ErrJSONFormat)
	}
	s, err := myjson.Decode[SendTo](d)
	if err != nil || s.To == 0 {
		log.Printf("Ошибка JSON")
		return s, fmt.Errorf("ParseSendTo: %w", err)
	}
	s.UserID = task.UserID
	s.ConnID = task.ConnID
	log.Println("отправить к ", s.To, "отправить от", s.UserID)
	return s, nil
}

func HandleSendTo(hubMap *HubMap, currentHub *Hub, t task.Task) error {
	s, err := ParseSendTo(t)
	if err != nil {
		return fmt.Errorf("HandleSendTo: %w", err)
	}

	h := hubMap.GetHub(s.To)
	h.mu.RLock()
	defer h.mu.RUnlock()
	slice := h.Map[s.To]
	log.Println(slice)
	for _, cl := range slice {
		err = cl.SafeSend([]byte(s.Message))
		switch err {
		case errs.ErrClientChanClosed:
		case errs.ErrClientChanFull:
			cl.CloseConn()
		case nil:
		}
	}
	return nil
}
