package repository

import "errors"

var (
	ErrNotFound   = errors.New("data not found")
	ErrEmailTaken = errors.New("email is already registered")
)
