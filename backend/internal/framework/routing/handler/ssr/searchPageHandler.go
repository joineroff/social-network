package ssr

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/framework/routing/middleware"
	"github.com/joineroff/social-network/backend/internal/usecase"
)

func SearchHandler(uc *usecase.SearchProfilesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		if !ok || userID.(string) == "" {
			c.Redirect(http.StatusFound, "/sign-in")
			return
		}

		input := dto.SearchProfileInputDto{}

		err := c.ShouldBind(&input)
		if err != nil {
			c.HTML(http.StatusOK, "search/index.tmpl", gin.H{
				"error": err.Error(),
			})

			return
		}

		offset := 0

		if input.Limit == 0 {
			input.Limit = 50
		}

		if input.Page < 1 {
			input.Page = 1
		}

		offset = (input.Page - 1) * input.Limit
		input.Offset = offset

		input.CurrentUserID = userID.(string)

		output, err := uc.Do(c.Request.Context(), &input)
		if err != nil {
			c.HTML(http.StatusOK, "search/index.tmpl", gin.H{
				"searchQuery": input.Query,
			})

			return
		}

		maxPage := (output.Total / input.Limit)
		if output.Total%input.Limit != 0 {
			maxPage += 1
		}

		lo, hi := 1, maxPage

		pages := make([]int, hi-lo+1)
		for i := range pages {
			pages[i] = i + lo
		}

		paginUri := fmt.Sprintf("/search?query=%s&limit=%d&page=", input.Query, input.Limit)
		profileUri := "/profile/?id="

		c.HTML(http.StatusOK, "search/index.tmpl", gin.H{
			"searchQuery":   input.Query,
			"limit":         input.Limit,
			"offset":        offset,
			"page":          input.Page,
			"pages":         pages,
			"paginationUri": paginUri,
			"profileUri":    profileUri,
			"profiles":      output.Profiles,
		})
	}
}

func SearchPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(middleware.UserIDKey)

		fmt.Println(userID)

		if !ok || userID.(string) == "" {
			c.Redirect(http.StatusFound, "/sign-in")
			return
		}

		c.HTML(http.StatusOK, "search/index.tmpl", gin.H{})
	}
}
