package utils

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal error")
	ErrNotModified  = errors.New("data not modified")
)

func MapError(err error) (int, string) {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound, err.Error()
	case ErrInvalidInput:
		return http.StatusBadRequest, err.Error()
	case ErrUnauthorized:
		return http.StatusUnauthorized, err.Error()
	case ErrConflict:
		return http.StatusConflict, err.Error()
	case ErrNotModified:
		return http.StatusNotModified, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
