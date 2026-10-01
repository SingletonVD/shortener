package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/SingletonVD/shortener/internal/apperror"
	"github.com/SingletonVD/shortener/internal/model"
	"github.com/google/uuid"
)

type userContextKeyT struct{}

var (
	userContextKey = userContextKeyT{}
)

const (
	sessionJWTCookie = "sessionJWT"
)

type TokenManager interface {
	Build(userID string) (string, error)
	Parse(tokenString string) (userID string, err error)
}

type AuthMiddleware struct {
	tokenManager TokenManager
}

func NewAuthMiddleware(tokenManager TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenManager}
}

func SetUserToContext(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	userUntyped := ctx.Value(userContextKey)
	if userUntyped == nil {
		return nil, apperror.ErrUserNotInContext
	}
	user, ok := userUntyped.(*model.User)
	if !ok {
		return nil, apperror.ErrUserIncorrectType
	}
	return user, nil
}

func (authMiddleware *AuthMiddleware) IntrospectUserJWT(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie(sessionJWTCookie)

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				authMiddleware.createUser(nextHandler).ServeHTTP(w, r)
				return
			}
			nextHandler.ServeHTTP(w, r)
			return
		}

		userID, err := authMiddleware.tokenManager.Parse(jwtCookie.Value)
		if err != nil {
			authMiddleware.createUser(nextHandler).ServeHTTP(w, r)
			return
		}

		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := model.User{UserID: userID}
		ctxWithUser := SetUserToContext(r.Context(), &user)
		nextHandler.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}

func (authMiddleware *AuthMiddleware) createUser(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newUserID := uuid.New().String()
		newUser := &model.User{UserID: newUserID}
		token, err := authMiddleware.tokenManager.Build(newUserID)

		if err != nil {
			nextHandler.ServeHTTP(w, r)
			return
		}

		ctxWithUser := SetUserToContext(r.Context(), newUser)

		http.SetCookie(w, &http.Cookie{
			Name:     string(sessionJWTCookie),
			Value:    token,
			HttpOnly: true,
			Path:     "/",
		})

		nextHandler.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}
