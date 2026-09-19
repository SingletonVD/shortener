package apperror

import (
	"errors"
	"fmt"
)

type FullLinkConflict struct {
	ShortLink string
	FullLink  string
}

func (conflict *FullLinkConflict) Error() string {
	return fmt.Sprintf("full link %s already exists as %s", conflict.FullLink, conflict.ShortLink)
}

type UnexpectedSigningMethod struct {
	Method string
}

func (unexpectedSigningMethod *UnexpectedSigningMethod) Error() string {
	return fmt.Sprintf("unexpected signing method: %s", unexpectedSigningMethod.Method)
}

var (
	TokenNotValid    = errors.New("token is not valid")
	UserIDNotDefined = errors.New("user id is not defined in token")
)
