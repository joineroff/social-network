package dto

type SignInInputDto struct {
	Login    string `json:"login" form:"login" validate:"required,min=3"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

type SignInOutputDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
