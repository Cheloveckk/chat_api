package socket

import "sync/atomic"

var ID int32 = 0

func GetNewID() int {
	return int(atomic.AddInt32(&ID, 1))
}
