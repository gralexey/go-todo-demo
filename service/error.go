package service

import "net/http"

type Error struct {
	suggestedCode int
	text          string
}

func (e Error) Error() string {
	return e.text
}

func NewErrorFromDBError(err error) error {
	return &Error{
		suggestedCode: http.StatusInternalServerError,
		text:          err.Error(),
	}
}

func NewErrorNotFound(err error) error {
	return &Error{
		suggestedCode: http.StatusNotFound,
		text:          "not found",
	}
}

func NewErrorAccessDenied() error {
	return &Error{
		suggestedCode: http.StatusForbidden,
		text:          "access denied",
	}
}
