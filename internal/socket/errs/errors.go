package errs

import "errors"

var ErrClientChanClosed = errors.New("client channel is closed")

var ErrClientChanFull = errors.New("client channel is full")

var ErrInvalidCommand = errors.New("command doesn`t exist")

var ErrNotauthorized = errors.New("not authorized")

var ErrInvalidChatID = errors.New("invalid chat id")

var ErrServerError = errors.New("server error")
