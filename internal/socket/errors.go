package socket

import "errors"

var ErrClientChanClosed = errors.New("client channel is closed")

var ErrClientChanFull = errors.New("client channel is full")

var ErrJSONFormat = errors.New("json request error")

var ErrInvalidCommand = errors.New("command doesn`t exist")

var ErrNotauthorized = errors.New("not authorized")
