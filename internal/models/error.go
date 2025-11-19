package models

import "fmt"

type Error struct {
	StatusCode  int
	Category    string
	Description string
}

type ErrorBucket struct {
	Err error
}

func NewError(statusCode int, category, description string) *Error {
	return &Error{
		StatusCode:  statusCode,
		Category:    category,
		Description: description,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Category, e.Description)
}
