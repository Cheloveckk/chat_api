package task

import (
	"encoding/json"
)

type Task struct {
	ConnID int64
	UserID int
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
}
