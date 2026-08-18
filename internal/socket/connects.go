package socket

//Unuse
import (
	"chat/api/internal/work"

	"github.com/gorilla/websocket"
)

type ShardedMap struct {
	Work   chan work.WorkType
	Count  int
	Shards []ConnsMap
}

type ConnsMap map[int]*User

type UsersMap map[int][]*websocket.Conn

func NewShardedMap(count int, len int) *ShardedMap {
	slice := make([]ConnsMap, 0, count)
	for i := 0; i < count; i++ {
		slice = append(slice, make(ConnsMap, len))
	}
	return &ShardedMap{
		Work:   make(chan work.WorkType, 5),
		Count:  count,
		Shards: slice,
	}
}

func (m *ShardedMap) Add(connId int, userId int, ws *websocket.Conn) {
	shard := connId % m.Count
	m.Shards[shard][connId] = &User{connId, userId, ws}
}

func (m *ShardedMap) Delete(id int) {
	shard := id % m.Count
	delete(m.Shards[shard], id)
}
