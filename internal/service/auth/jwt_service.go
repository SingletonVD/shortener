package auth

import (
	"fmt"

	"github.com/SingletonVD/shortener/internal/apperror"
	"github.com/SingletonVD/shortener/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	authSecret []byte
}

type SigningMethod = jwt.SigningMethodHMAC

var (
	signingMethod = jwt.SigningMethodHS256
)

func NewAuthService(authSecret string) *AuthService {
	return &AuthService{authSecret: []byte(authSecret)}
}

func (authService *AuthService) IntrospectToken(tokenString string) (*model.User, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*SigningMethod); !ok {
			return nil, &apperror.UnexpectedSigningMethod{Method: fmt.Sprintf("%v", t.Header["alg"])}
		}
		return authService.authSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, apperror.TokenNotValid
	}

	if claims.Subject == "" {
		return nil, apperror.UserIDNotDefined
	}
	return &model.User{UserID: claims.Subject}, nil
}

func (authService *AuthService) CreateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(signingMethod, jwt.RegisteredClaims{
		Subject: user.UserID,
	})

	tokenString, err := token.SignedString(authService.authSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
