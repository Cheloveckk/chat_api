package main

import (
	"chat/api/config"
	"chat/api/internal/db"
	"chat/api/internal/socket"
	"chat/api/pkg/hello"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/ws", handleRequest)
	hello.HelloHandle(handler)
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}
	conf := config.GetConfig()
	fmt.Println(conf.DbConfig.Key)
	dbConn := db.GetDbConn(conf)
	err := dbConn.Ping()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Succes")
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ws, err := socket.InitSocket(w, r)
	if err != nil {
		log.Fatal(err.Error())
	}
	ws.WriteJSON("message")
	go func(ws *websocket.Conn) {
		var s string
		for {
			fmt.Scan(&s)
			ws.WriteJSON(s)
		}
	}(ws)
	go func(ws *websocket.Conn) {
		data := map[string]any{}
		for {
			err := ws.ReadJSON(&data)
			if _, ok := err.(*websocket.CloseError); ok == true {
				ws.Close()
				break
			} else if err == nil {
				fmt.Println(data)
			}
		}
	}(ws)
}
