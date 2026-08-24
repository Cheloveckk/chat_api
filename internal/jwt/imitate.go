package jwt

import (
	"net/http"
	"strconv"
	"sync/atomic"
)

var ID int64 = 0

func GetID_(r *http.Request) int {
	return int(atomic.AddInt64(&ID, 1))
}

func GetID(r *http.Request) int {
	id, err := strconv.Atoi(r.Header.Get("Authorization"))
	if err != nil {
		return 0
	}
	return id
}
