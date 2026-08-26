package myjson

import (
	"encoding/json"
	"errors"
)

var ErrJSONFormat = errors.New("json parse error")

func Decode[T any](data []byte) (T, error) {
	var value T
	err := json.Unmarshal(data, &value)
	if err != nil {
		return value, ErrJSONFormat
	}
	return value, err
}
