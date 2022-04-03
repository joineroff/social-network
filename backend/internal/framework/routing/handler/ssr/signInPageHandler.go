package ssr

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/framework/routing/middleware"
	"github.com/joineroff/social-network/backend/internal/usecase"
)

func SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		if !ok || userID.(string) == "" {
			c.HTML(http.StatusOK, "sign-in/index.tmpl", gin.H{})
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}

func SignInPostHandler(domain string, uc *usecase.SignInUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &dto.SignInInputDto{}

		err := c.ShouldBind(&input)
		if err != nil {
			c.HTML(http.StatusOK, "sign-in/index.tmpl", gin.H{
				"error": err.Error(),
			})

			return
		}

		output, err := uc.Do(c.Request.Context(), input)
		if err != nil {
			c.HTML(http.StatusOK, "sign-in/index.tmpl", gin.H{
				"login": input.Login,
				"error": err.Error(),
			})

			return
		}

		maxCookieTTL := 3600
		c.SetCookie(middleware.CookieTokenKey, output.AccessToken, maxCookieTTL, "/", domain, false, true)
		c.Redirect(http.StatusFound, "/")
	}
}
