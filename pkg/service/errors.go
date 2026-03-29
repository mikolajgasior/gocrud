package service

import (
	"errors"
)

const (
	LogAttrService = "service"
	LogAttrError   = "err"
)

var NotFoundError = errors.New("item not found")
var InvalidFilterOpError = errors.New("invalid filter op")

type ModelValidationError struct {
	Err        error
	Violations map[string]uint64
}

func (e *ModelValidationError) Error() string {
	return e.Err.Error()
}

type FilterValidationError struct {
	Err        error
	Violations map[string]uint64
}

func (e *FilterValidationError) Error() string {
	return e.Err.Error()
}
