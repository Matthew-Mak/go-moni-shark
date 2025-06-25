package errors

import "errors"

var (
	ErrCreateFile      = errors.New("failed to create file")
	ErrWriteLine       = errors.New("failed to write line")
	ErrInvolvingWriter = errors.New("failed involving writer")
	ErrOpenFile        = errors.New("failed to open file")
	ErrConvertString   = errors.New("failed convert string")
	ErrParseBool       = errors.New("failed parse bool")
	ErrParseTime       = errors.New("failed parse time")
)
