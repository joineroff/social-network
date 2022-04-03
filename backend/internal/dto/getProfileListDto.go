package dto

type GetProfileListInputDto struct {
	CurrentUserID string `json:"-"`
	Limit         int64  `json:"limit"`
	Offset        int64  `json:"offset"`
}

type GetProfileListOutputDto = []ProfileDto
