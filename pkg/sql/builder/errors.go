package builder

type BuilderError struct {
	Op  string
	Tag string
	Err error
}

func (e *BuilderError) Error() string {
	return e.Op + ": " + e.Err.Error()
}

var getQueryBuilderError = func(op string, err error) *BuilderError {
	return &BuilderError{
		Op:  "QueryBuilder." + op,
		Err: err,
	}
}
