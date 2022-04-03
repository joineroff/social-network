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

func HomeHandler(domain string, uc *usecase.GetProfileUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		fmt.Println(userID)

		if !ok || userID.(string) == "" {
			c.Redirect(http.StatusFound, "/sign-in")
			return
		}

		input := &dto.GetProfileInputDto{
			ID: userID.(string),
		}

		profile, err := uc.Do(c.Request.Context(), input)
		if err != nil {
			c.SetCookie(middleware.CookieTokenKey, "", 0, "/", domain, false, true)
			c.Redirect(http.StatusFound, "/sign-in")

			return
		}

		c.HTML(http.StatusOK, "home/index.tmpl", gin.H{
			"login":     profile.Profile.User.Login,
			"firstName": profile.Profile.User.FirstName,
			"lastName":  profile.Profile.User.LastName,
			"gender":    profile.Profile.User.GenderAsStirng(),
			"interests": strings.Join(profile.Profile.User.Interests, ", "),
			"city":      profile.Profile.User.City,
		})
	}
}
