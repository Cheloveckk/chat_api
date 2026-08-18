package work

import "encoding/json"

type WorkType struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
