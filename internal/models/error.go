package models

import "fmt"

type Error struct {
	StatusCode  int
	Category    string
	Descritpion string
}

type ErrorBucket struct {
	Err error
}

func NewError(statusCode int, category string, description string) *Error {
	return &Error{
		StatusCode:  statusCode,
		Category:    category,
		Descritpion: description,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Category, e.Descritpion)
}