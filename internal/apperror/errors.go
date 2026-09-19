package apperror

import (
	"errors"
	"fmt"
)

type ErrFullLinkConflict struct {
	ShortLink string
	FullLink  string
}

func (conflict *ErrFullLinkConflict) Error() string {
	return fmt.Sprintf("full link %s already exists as %s", conflict.FullLink, conflict.ShortLink)
}

type ErrUnexpectedSigningMethod struct {
	Method string
}

func (unexpectedSigningMethod *ErrUnexpectedSigningMethod) Error() string {
	return fmt.Sprintf("unexpected signing method: %s", unexpectedSigningMethod.Method)
}

var (
	ErrTokenNotValid    = errors.New("token is not valid")
	ErrUserIDNotDefined = errors.New("user id is not defined in token")
)
