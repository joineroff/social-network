package ssr

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/framework/routing/middleware"
	"github.com/joineroff/social-network/backend/internal/usecase"
)

func ProfileHandler(uc *usecase.GetProfileUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		if !ok || userID.(string) == "" {
			c.Redirect(http.StatusFound, "/sign-in")
			return
		}

		profileID := c.Param("profileID")

		input := &dto.GetProfileInputDto{
			ID:            profileID,
			CurrentUserID: userID.(string),
		}

		profile, err := uc.Do(c.Request.Context(), input)
		if err != nil {
			c.Redirect(http.StatusFound, "/")

			return
		}

		c.HTML(http.StatusOK, "profile/index.tmpl", gin.H{
			"login":     profile.Profile.User.Login,
			"firstName": profile.Profile.User.FirstName,
			"lastName":  profile.Profile.User.LastName,
			"gender":    profile.Profile.User.GenderAsStirng(),
			"interests": strings.Join(profile.Profile.User.Interests, ", "),
			"city":      profile.Profile.User.City,
			"isFriend":  profile.Profile.IsFriend,
			"id":        profile.Profile.User.ID,
		})
	}
}

func ProfileAddFriendHandler(uc *usecase.AddFriendUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		if !ok || userID.(string) == "" {
			c.Redirect(http.StatusFound, "/sign-in")
			return
		}

		profileID := c.Param("profileID")

		input := &dto.AddFriendInputDto{
			FriendID:      profileID,
			CurrentUserID: userID.(string),
		}

		if err := uc.Do(c.Request.Context(), input); err != nil {
			c.Redirect(http.StatusFound, "/profile/"+profileID)

			return
		}

		c.Redirect(http.StatusFound, "/profile/"+profileID)
	}
}

func ProfileRemoveFriendHandler(uc *usecase.RemoveFriendUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		if !ok || userID.(string) == "" {
			c.Redirect(http.StatusFound, "/sign-in")
			return
		}

		profileID := c.Param("profileID")

		input := &dto.RemoveFriendInputDto{
			FriendID:      profileID,
			CurrentUserID: userID.(string),
		}

		if err := uc.Do(c.Request.Context(), input); err != nil {
			c.Redirect(http.StatusFound, "/profile/"+profileID)

			return
		}

		c.Redirect(http.StatusFound, "/profile/"+profileID)
	}
}
