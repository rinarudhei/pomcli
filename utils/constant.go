// Package utils contain custom error used in Pomcli App
package utils

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrSessionNotFound = errors.New("no session found")
)
