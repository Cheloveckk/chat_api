package db

import (
	"chat/api/config"
	"database/sql"

	_ "github.com/lib/pq"
)

type Db struct {
	*sql.DB
}

func GetDbConn(conf *config.Config) *Db {
	con, err := sql.Open("postgres", conf.DbConfig.Key)
	if err != nil {
		panic(err.Error())
	}
	return &Db{DB: con}
}
