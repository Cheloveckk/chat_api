package main

import (
	"chat/api/config"
	"chat/api/pkgs/mydb"
	"fmt"
)

func main() {
	conf := config.GetConfig()
	db := mydb.GetDbConn(conf)
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`DROP TABLE IF EXISTS Users`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Успех")
}
