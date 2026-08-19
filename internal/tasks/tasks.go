package tasks

import "encoding/json"

type Task struct {
	From int64
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
