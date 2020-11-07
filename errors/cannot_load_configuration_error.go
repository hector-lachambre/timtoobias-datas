package errors

type CannotLoadConfigurationError struct {
	Message string
}

func (e *CannotLoadConfigurationError) Error() string {
	return e.Message
}
