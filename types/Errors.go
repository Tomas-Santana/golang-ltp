package types

import "errors"

type Errors error

var ErrInvalidLogLevel = errors.New("invalid log level")
var ErrInvalidResponseStatus = errors.New("invalid response status")
var ErrRequestTooLong = errors.New("request too long")
var ErrInvalidRequest = errors.New("invalid request")


