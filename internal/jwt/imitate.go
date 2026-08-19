package jwt

import (
	"net/http"
	"sync/atomic"
)

var ID int64 = 0

func GetID(r *http.Request) int64 {
	return atomic.AddInt64(&ID, 1)
}
