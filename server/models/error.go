package models

// Error model
type Error struct {
	Message string
}

// NewError creates an instance
func NewError(message string) *Error {
	return &Error{Message: message}
}

func (e *Error) Error() string {
	return e.Message
}
