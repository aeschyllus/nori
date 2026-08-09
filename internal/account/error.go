package account

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInvalidAccount  = errors.New("invalid account")
)
