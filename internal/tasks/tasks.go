package tasks

import "encoding/json"

type Task struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
