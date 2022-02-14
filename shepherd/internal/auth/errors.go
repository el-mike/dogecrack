package auth

import "fmt"

// UserNotFoundError - thrown when given User does not exist.
type UserNotFoundError struct {
	Name string `json:"name"`
}

// UserNotFoundError - returns a new UserNotFoundError instance.
func NewUserNotFoundError(name string) *UserNotFoundError {
	return &UserNotFoundError{
		Name: name,
	}
}

// Error - Error interface implementation.
func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("User '%s' does not exist", e.Name)
}

// IncorrectPasswordError - thrown when User provides incorrect password.
type IncorrectPasswordError struct{}

// IncorrectPasswordError - returns a new IncorrectPasswordError instance.
func NewIncorrectPasswordError() *IncorrectPasswordError {
	return &IncorrectPasswordError{}
}

// Error - Error interface implementation.
func (e *IncorrectPasswordError) Error() string {
	return fmt.Sprintf("Provided password is incorrect")
}

// UserUnathorizedError - thrown when a client is not authenticated / authorized.
type UserUnathorizedError struct{}

// UserUnathorizedError - returns a new UserUnathorizedError instance.
func NewUserUnathorizedError() *UserUnathorizedError {
	return &UserUnathorizedError{}
}

// Error - Error interface implementation.
func (e *UserUnathorizedError) Error() string {
	return fmt.Sprintf("You are not authorized")
}

// SessionExpiredError - thrown when user session expired.
type SessionExpiredError struct{}

// SessionExpiredError - returns a new SessionExpiredError instance.
func NewSessionExpiredError() *SessionExpiredError {
	return &SessionExpiredError{}
}

// Error - Error interface implementation.
func (e *SessionExpiredError) Error() string {
	return fmt.Sprintf("Your session expired, log in again")
}

// SessionCookieMissingError - thrown when request lacks a session cookie.
type SessionCookieMissingError struct{}

// SessionCookieMissingError - returns a new SessionCookieMissingError instance.
func NewSessionCookieMissingError() *SessionCookieMissingError {
	return &SessionCookieMissingError{}
}

// Error - Error interface implementation.
func (e *SessionCookieMissingError) Error() string {
	return fmt.Sprintf("Session cookie is missing!")
}
