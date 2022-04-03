package usecase

import (
	"context"

	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/service"
	"github.com/joineroff/social-network/backend/internal/validation"
)

type GetProfileUsecase struct {
	profileService service.ProfileService
}

func NewGetProfileUsecase(
	profileService service.ProfileService,
) *GetProfileUsecase {
	return &GetProfileUsecase{
		profileService: profileService,
	}
}

func (uc *GetProfileUsecase) Do(ctx context.Context, input *dto.GetProfileInputDto) (*dto.GetProfileOutputDto, error) {
	output := &dto.GetProfileOutputDto{}

	if err := validation.Struct(input); err != nil {
		return nil, err
	}

	profile, err := uc.profileService.GetProfile(ctx, input.ID, input.CurrentUserID)
	if err != nil {
		return nil, err
	}

	output.Profile = profile

	return output, nil
}
