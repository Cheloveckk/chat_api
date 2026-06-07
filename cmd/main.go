package main

import (
	"chat/api/pkg/hello"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var Upg = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/ws")
	hello.HelloHandle(handler)
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ws, err := Upg.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}
