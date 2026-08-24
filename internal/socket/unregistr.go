package socket

import "log"

type unregData struct {
	connID int
	userID int
}

func unregWorker(h *Hub) {
	defer h.wg.Done()
	for {
		select {
		case data, ok := <-h.Unreg:
			if !ok {
				return
			}
			log.Println("unreg worker data:", data)
			log.Println("Попытка удалить пользователя")
			h.DeleteConn(data.connID, data.userID)
		case <-h.Ctx.Done():
			return
		}
	}
}
