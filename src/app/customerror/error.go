package customerror

import "github.com/pkg/errors"

var (
	ErrAlreadyExists  = errors.New("record already exists")
	ErrInvalidInput   = errors.New("invalid input")
	ErrFileNotFound   = errors.New("file not found")
	ErrNoDataProvided = errors.New("no data provided")
	ErrNoSubscribers  = errors.New("no subscribers in storage")
)
