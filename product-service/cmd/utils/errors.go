package utils

// Make error struct to implement error interface
type Error struct {
	Message string
}

// Implement Error() method on the Error struct
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new error
func NewError(message string) *Error {
	msg := message
	return &Error{Message: msg}
}
