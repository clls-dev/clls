package lspsrv

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrUnknownMethod  = errors.New("unknown method")
)
