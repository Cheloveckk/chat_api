package socket

import (
	"chat/api/internal/task"
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
	var s SendTo
	err := json.Unmarshal(task.Data, &s)
	if err != nil || s.To == 0 {
		log.Printf("Ошибка JSON")
		return s, ErrJSONFormat
	}
	s.UserID = task.UserID
	s.ConnID = task.ConnID
	fmt.Println("отправить к ", s.To, "отправить от", s.UserID)
	return s, nil
}

func HandleSendTo(hubMap *HubMap, currentHub *Hub, t task.Task) error {
	s, err := ParseSendTo(t)
	if err != nil {
		return err
	}

	h := hubMap.GetHub(s.To)
	h.mu.RLock()
	defer h.mu.RUnlock()
	slice := h.Map[s.To]
	fmt.Println(slice)
	for _, cl := range slice {
		err = cl.SafeSend([]byte(s.Message))
		switch err {
		case ErrClientChanClosed:
			return nil
		case ErrClientChanFull:
			go func() {
				select {
				case h.Unreg <- unregData{connID: cl.ConnID, userID: cl.UserID}:
					log.Println("Пользователь неактивен")
				case <-h.Ctx.Done():
					log.Println("Контекст не текущего хаба отменен")
				case <-currentHub.Ctx.Done():
					log.Println("Контекст текущего хаба отменен")
				}
			}()
			return err
		case nil:
		}
	}

	return nil
}
