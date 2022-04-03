package ssr

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joineroff/social-network/backend/internal/framework/routing/middleware"
)

func SignOutGetHandler(domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		if !ok || userID.(string) == "" {
			c.Redirect(http.StatusFound, "/sign-in")
			return
		}

		c.SetCookie(middleware.CookieTokenKey, "", 0, "/", domain, false, true)

		c.Redirect(http.StatusFound, "/sign-in")
	}
}

func SignOutPostHandler(domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		if !ok || userID.(string) == "" {
			c.Redirect(http.StatusFound, "/sign-in")
			return
		}

		c.SetCookie(middleware.CookieTokenKey, "", 0, "/", domain, false, true)

		c.Redirect(http.StatusFound, "/sign-in")
	}
}
