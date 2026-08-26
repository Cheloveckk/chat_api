package socket

import (
	"chat/api/internal/socket/errs"
	"chat/api/internal/task"
	"chat/api/pkgs/myjson"
	"encoding/json"
	"fmt"
	"log"
)

type SendChat struct {
	ChatID  int    `json:"chat"`
	Message string `json:"message"`
	From    int    `json:"from"`
}

type DeliverySendChat struct {
	Send  SendChat
	Users []int
}

func ParseSendChat(data any) (SendChat, error) {
	d, ok := data.(json.RawMessage)
	if !ok {
		return SendChat{}, fmt.Errorf("ParseSendTo: %w", errs.ErrServerError)
	}
	s, err := myjson.Decode[SendChat](d)
	if err != nil {
		return s, fmt.Errorf("ParseSendTo: %w", err)
	}
	return s, nil
}

// func ParseDeliverySendChat(data any) (DeliverySendChat, error) {
// 	d, ok := data.([]byte)
// 	if !ok {
// 		return DeliverySendChat{}, fmt.Errorf("DeliveryParseSendTo: %w", errs.ErrServerError)
// 	}
// 	s, err := myjson.Decode[DeliverySendChat](d)
// 	if err != nil {
// 		return s, fmt.Errorf("DeliveryParseSendTo: %w", err)
// 	}
// 	return s, nil
// }

// ИЗМЕНИИИИИИТЬ
func HandleSendChat(hubMap *HubMap, currentHub *Hub, t task.Task) error {
	log.Println("HandleSendChat")
	s, err := ParseSendChat(t.Data)
	if err != nil {
		log.Println("err")
		return fmt.Errorf("HandleSendChat: %w", err)
	}

	if s.ChatID < 0 {
		return errs.ErrInvalidChatID
	}
	s.From = t.UserID
	if s.ChatID == 0 {
		hubMap.mu.RLock()
		for _, h := range hubMap.Hubs {
			h := h
			go func() {
				log.Println("Delivery")
				h.Tasks <- task.Task{
					Type: "delivery",
					Data: DeliverySendChat{
						s,
						nil,
					},
				}
			}()
		}
		hubMap.mu.RUnlock()
	} else {
	}
	return nil
}
func HandleDeliverySendChat(h *Hub, data any) {
	delivery, ok := data.(DeliverySendChat)
	if !ok {
		log.Fatal(fmt.Errorf("HandleDeliverySendChat: %w", errs.ErrServerError))
	}
	send, err := json.Marshal(delivery.Send)
	if err != nil {
		log.Fatal(fmt.Errorf("HandleDeliverySendChat: %w", myjson.ErrJSONFormat))
	}
	if delivery.Send.ChatID == 0 {
		h.mu.RLock()
		for _, userSlice := range h.Map {
			for _, client := range userSlice {
				err := client.SafeSend(send)
				if err != nil {
					client.CloseConn()
				}
			}
		}
		h.mu.RUnlock()
		return
	}
	return
}
