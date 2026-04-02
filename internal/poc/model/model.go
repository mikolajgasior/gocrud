package model

type ValidationError struct {
	Violations map[string]uint64
}

func (e ValidationError) Error() string {
	return "validation error: TODO"
}
