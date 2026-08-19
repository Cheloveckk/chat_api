package main

import (
	"chat/api/config"
	"chat/api/internal/db"
	"chat/api/internal/socket"
	"chat/api/pkg/hello"
	"fmt"
	"log"
	"net/http"
)

var hubMap = socket.NewHubMap(5, 5)

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
	cl, err := socket.InitConn(w, r)
	if err != nil {
		log.Fatal(err.Error())
	}
	cl.Conn.WriteJSON("message")
	cl.StartWork(hubMap.GetHub(cl.UserID).ChTasks)
	hubMap.AddClient(cl)
}
