package constants

type ErrorMessage string

const (
	ErrNotFound            ErrorMessage = "Resource not found"
	ErrInternalServerError ErrorMessage = "Something went wrong :<"
)
