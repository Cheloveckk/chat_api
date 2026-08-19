package socket

import "sync/atomic"

var ID int64 = 0

func GetNewID() int64 {
	return atomic.AddInt64(&ID, 1)
}
