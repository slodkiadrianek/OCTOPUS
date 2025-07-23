package Models

type Error struct {
	StatusCode  int
	Category    string
	Descritpion string
}

func NewError(statusCode int, category string, description string) *Error {
	return &Error{
		StatusCode:  statusCode,
		Category:    category,
		Descritpion: description,
	}
}
