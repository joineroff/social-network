package dto

import "github.com/joineroff/social-network/backend/internal/entity"

type GetProfileInputDto struct {
	ID            string `json:"id"`
	CurrentUserID string `json:"-"`
}

type GetProfileOutputDto struct {
	Profile *entity.Profile `json:"profile"`
}
