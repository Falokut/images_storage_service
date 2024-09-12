package domain

import (
	"errors"
)

var (
	ErrImageNotFound = errors.New("image not found")
)

const (
	ErrCodeImageNotFound       = 600
	ErrCodeImageTooBig         = 601
	ErrCodeImageHasZeroSize    = 602
	ErrCodeUnsupportedFileType = 603
)

type InvalidArgumentError struct {
	ErrCode int
	Reason  string
}

func NewInvalidArgumentError(reason string, errCode int) InvalidArgumentError {
	return InvalidArgumentError{Reason: reason, ErrCode: errCode}
}

func (e InvalidArgumentError) Error() string {
	return e.Reason
}
