package auth

import "fmt"

type UserNotFoundError struct {
	Name string `json:"name"`
}

func NewUserNotFoundError(name string) *UserNotFoundError {
	return &UserNotFoundError{
		Name: name,
	}
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("User '%s' does not exist", e.Name)
}

type IncorrectPasswordError struct{}

func NewIncorrectPasswordError() *IncorrectPasswordError {
	return &IncorrectPasswordError{}
}

func (e *IncorrectPasswordError) Error() string {
	return fmt.Sprintf("Provided password is incorrect")
}

type UserUnathorizedError struct{}

func NewUserUnathorizedError() *UserUnathorizedError {
	return &UserUnathorizedError{}
}

func (e *UserUnathorizedError) Error() string {
	return fmt.Sprintf("You are not authorized")
}

type SessionExpiredError struct{}

func NewSessionExpiredError() *SessionExpiredError {
	return &SessionExpiredError{}
}

func (e *SessionExpiredError) Error() string {
	return fmt.Sprintf("Your session expired, log in again")
}
