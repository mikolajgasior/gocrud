package crud

import (
	"errors"
	"fmt"
)

type ErrCRUD struct {
	Op  string
	Tag string
	Err error
}

func (e ErrCRUD) Error() string {
	return e.Err.Error()
}

func (e ErrCRUD) Unwrap() error {
	return e.Err
}

type ErrValidation struct {
	Violations map[string]int
	Err        error
}

func (e ErrValidation) Error() string {
	return e.Err.Error()
}

func (e ErrValidation) Unwrap() error {
	return e.Err
}

type ErrUniq struct {
	Field string
}

func (e ErrUniq) Error() string {
	return errors.New(fmt.Sprintf("uniq %s failed", e.Field)).Error()
}

func (e ErrUniq) Unwrap() error {
	return errors.New(fmt.Sprintf("uniq %s failed", e.Field))
}

func (e ErrUniq) Is(target error) bool {
	_, ok := target.(ErrUniq)
	return ok
}
