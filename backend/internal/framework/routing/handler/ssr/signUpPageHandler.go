package ssr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/framework/routing/middleware"
	"github.com/joineroff/social-network/backend/internal/usecase"
)

func SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		if !ok || userID.(string) == "" {
			c.HTML(http.StatusOK, "sign-up/index.tmpl", gin.H{})
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}

func SignUpPostHandler(domain string, uc *usecase.SignUpUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &dto.SignUpInputDto{}

		err := c.ShouldBind(&input)
		if err != nil {
			c.HTML(http.StatusOK, "sign-up/index.tmpl", gin.H{
				"error": err.Error(),
			})

			return
		}

		fmt.Print(input)

		output, err := uc.Do(c.Request.Context(), input)
		if err != nil {
			c.HTML(http.StatusOK, "sign-up/index.tmpl", gin.H{
				"login":     input.Login,
				"firstName": input.FirstName,
				"lastName":  input.LastName,
				"gender":    input.Gender,
				"interests": strings.Join(input.Interests, ", "),
				"city":      input.City,
				"error":     err.Error(),
			})

			return
		}

		maxCookieTTL := 3600
		c.SetCookie("token", output.AccessToken, maxCookieTTL, "/", domain, false, true)
		c.Redirect(http.StatusFound, "/")
	}
}
