package main

import (
	"chat/api/config"
	"chat/api/internal/mydb"
)

func main() {
	conf := config.GetConfig()

	db := mydb.GetDbConn(conf)
	err := db.Ping()

	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Exec(`CREATE TABLE IF NOT EXISTS Users(
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE,
		password VARCHAR(30) NOT NULL,
		user_name VARCHAR(30) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS Chats(
		id SERIAL PRIMARY KEY,
		name VARCHAR(30) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS Members(
		chat_id INTEGER,
		user_id INTEGER,
		role SMALLINT,
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		left_at TIMESTAMP DEFAULT NULL,
		FOREIGN KEY (chat_id) REFERENCES Chats (id)
		ON DELETE SET NULL ON UPDATE CASCADE,
		FOREIGN KEY (user_id) REFERENCES Users (id)
		ON DELETE SET NULL ON UPDATE CASCADE
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS Messages(
		chat_id INTEGER,
		user_id INTEGER,
		text VARCHAR(150),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL,
		FOREIGN KEY (chat_id) REFERENCES Chats (id)
		ON DELETE SET NULL ON UPDATE CASCADE,
		FOREIGN KEY (user_id) REFERENCES Users (id)
		ON DELETE SET NULL ON UPDATE CASCADE
	)`)
}
