package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joineroff/social-network/backend/internal/service"
)

const (
	UserIDKey      = "userIDAuthKey"
	CookieTokenKey = "token"
	HeaderTokenKey = "token"
)

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(
	authService service.AuthService,
) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (am *AuthMiddleware) ExtractUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			userID string
			err    error
		)

		userID, err = am.getUserFromHeaderAccessToken(c)
		if userID == "" || err != nil {
			userID, err = am.getUserFromCookieAccessToken(c)
		}

		if err != nil {
			return
		}

		if userID != "" {
			c.Set(UserIDKey, userID)
		}
	}
}

func (am *AuthMiddleware) getUserFromCookieAccessToken(c *gin.Context) (string, error) {
	token, err := c.Cookie(CookieTokenKey)
	if err != nil {
		return "", err
	}

	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	claims, err := am.authService.ParseToken(token)
	if err != nil {
		return "", err
	}

	return claims.ID, nil
}

func (am *AuthMiddleware) getUserFromHeaderAccessToken(c *gin.Context) (string, error) {
	token := c.Request.Header.Get(HeaderTokenKey)
	if token == "" {
		return "", fmt.Errorf("header token is not set")
	}

	claims, err := am.authService.ParseToken(token)
	if err != nil {
		return "", err
	}

	return claims.ID, nil
}
