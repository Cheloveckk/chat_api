package main

import (
	"chat/api/config"
	"chat/api/internal/socket"
	"chat/api/pkgs/hello"
	"log"
	"net/http"
	_ "net/http/pprof"
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
	log.Println(conf.DbConfig.Key)

	// dbConn := db.GetDbConn(conf)
	// err := dbConn.Ping()
	// if err != nil {
	// 	panic(err.Error())
	// } else {
	//  log.Println("Succes")
	// }
	hubMap.CreateWorkers(5)

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	// go func() {
	// 	time.Sleep(time.Second * 10)
	// 	log.Println("TIME OUT")
	// 	hubMap.StopWork()
	// }()
	err := server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	cl, err := socket.InitConn(w, r)
	if err != nil {
		return
	}
	cl.Conn.WriteJSON("message")
	cl.StartWork(hubMap.GetHub(cl.UserID))
	hubMap.AddClient(cl)
	log.Println(cl.ConnID)
}
