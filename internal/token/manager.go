package token

import (
	"fmt"

	"github.com/SingletonVD/shortener/internal/apperror"
	"github.com/golang-jwt/jwt/v4"
)

type Manager struct {
	secret        []byte
	signingMethod jwt.SigningMethod
}

func NewManager(secret string) *Manager {
	return &Manager{
		secret:        []byte(secret),
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (manager *Manager) Build(userID string) (string, error) {
	token := jwt.NewWithClaims(manager.signingMethod, jwt.RegisteredClaims{
		Subject: userID,
	})

	tokenString, err := token.SignedString(manager.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (manager *Manager) Parse(tokenString string) (userID string, err error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != manager.signingMethod.Alg() {
			return nil, &apperror.ErrUnexpectedSigningMethod{Method: fmt.Sprintf("%v", t.Header["alg"])}
		}
		return manager.secret, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", apperror.ErrTokenNotValid
	}

	return claims.Subject, nil
}
