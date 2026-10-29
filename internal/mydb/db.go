package mydb

import (
	"chat/api/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Db struct {
	*sqlx.DB
}

func GetDbConn(conf *config.Config) *Db {
	db, err := sqlx.Open("postgres", conf.DbConfig.Key)
	if err != nil {
		panic(err.Error())
	}
	return &Db{DB: db}
}
