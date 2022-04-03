package usecase

type Error struct {
	Message string
	Details map[string]string
}

func (e Error) Error() string {
	return e.Message
}

type InternalError struct {
	Message string
	Details map[string]string
}

type LogicError struct {
	Message string
	Details map[string]string
}

type ValidationError struct {
	Message string
	Details map[string]string
}

func NewLogicError(message string) LogicError {
	e := LogicError{
		Message: message,
		Details: make(map[string]string),
	}

	return e
}

func (e LogicError) Error() string {
	return e.Message
}

func NewInternalError(message string) InternalError {
	e := InternalError{
		Message: message,
		Details: make(map[string]string),
	}

	return e
}

func (e InternalError) Error() string {
	return e.Message
}

func NewValidationError(message string) ValidationError {
	e := ValidationError{
		Message: message,
		Details: make(map[string]string),
	}

	return e
}

func (e ValidationError) Error() string {
	return e.Message
}
