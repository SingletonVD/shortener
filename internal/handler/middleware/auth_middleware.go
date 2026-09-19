package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/SingletonVD/shortener/internal/apperror"
	"github.com/SingletonVD/shortener/internal/model"
	"github.com/SingletonVD/shortener/internal/service/auth"
	"github.com/google/uuid"
)

type UserContextKey struct{}

var (
	userContextKey = UserContextKey{}
)

const (
	sessionJWTCookie = "sessionJWT"
)

type AuthMiddleware struct {
	authService *auth.AuthService
}

func NewAuthMiddleware(authService *auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func SetUserToContext(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)
	return user, ok
}

func (authMiddleware *AuthMiddleware) IntrospectUserJWT(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie(sessionJWTCookie)

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				authMiddleware.createUser(nextHandler).ServeHTTP(w, r)
			}
			nextHandler.ServeHTTP(w, r)
			return
		}

		user, err := authMiddleware.authService.IntrospectToken(jwtCookie.Value)
		if err != nil {
			if errors.Is(err, apperror.ErrUserIDNotDefined) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			nextHandler.ServeHTTP(w, r)
			return
		}

		ctxWithUser := SetUserToContext(r.Context(), user)
		nextHandler.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}

func (authMiddleware *AuthMiddleware) createUser(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newUserID := uuid.New().String()
		newUser := &model.User{UserID: newUserID}
		token, err := authMiddleware.authService.CreateToken(newUser)

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
