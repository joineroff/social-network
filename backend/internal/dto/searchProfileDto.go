package dto

import "github.com/joineroff/social-network/backend/internal/entity"

type SearchProfileInputDto struct {
	Query         string `json:"q" form:"q"`
	OnlyFriends   bool   `json:"onlyFriends" form:"onlyFriends"`
	CurrentUserID string `json:"-"`

	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit" validate:"max=100"`

	Page int `form:"page, omitempty"`
}

type SearchProfileOutputDto struct {
	Profiles []*entity.Profile `json:"profiles"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`

	Total int `json:"total"`
}
