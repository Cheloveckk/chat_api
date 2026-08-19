package send

import "encoding/json"

type SendTo struct {
	To      int    `json:"to"`
	From    int    `json:"from"`
	Message string `json:"message"`
}

func ParseSendTo(data []byte) (*SendTo, error) {
	var s SendTo
	err := json.Unmarshal(data, &s)
	return &s, err
}
