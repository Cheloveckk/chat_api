package task

import "encoding/json"

type Task struct {
	ConnID int
	UserID int
	Type   string `json:"type"`
	Data   any    `json:"data"`
}

type ReaderTask struct {
	ConnID int
	UserID int
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
}
