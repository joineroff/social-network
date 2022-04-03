package dto

type SignUpInputDto struct {
	Login     string   `json:"login" form:"login" validate:"required,min=3"`
	Password  string   `json:"password" form:"password" validate:"required,min=8"`
	FirstName string   `json:"firstName" form:"firstName" validate:"required,min=3"`
	LastName  string   `json:"lastName" form:"lastName"`
	Gender    int8     `json:"gender" form:"gender" validate:"min=0,max=8"`
	Interests []string `json:"interests" form:"interests"`
	City      string   `json:"city" form:"city" validate:"required,min=2"`
}

type SignUpOutputDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
