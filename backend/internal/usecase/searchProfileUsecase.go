package usecase

import (
	"context"

	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/service"
	"github.com/joineroff/social-network/backend/internal/validation"
	"go.uber.org/zap"
)

const minSearchProfileLimit = 10

type SearchProfilesUsecase struct {
	profileService service.ProfileService
}

func NewSearchProfilesUsecase(
	profileService service.ProfileService,
) *SearchProfilesUsecase {
	return &SearchProfilesUsecase{
		profileService: profileService,
	}
}

func (uc *SearchProfilesUsecase) Do(ctx context.Context, input *dto.SearchProfileInputDto) (*dto.SearchProfileOutputDto, error) {
	z := zap.S().With("context", "SearchProfilesUsecase")

	if err := validation.Struct(input); err != nil {
		return nil, err
	}

	if input.Limit < minSearchProfileLimit {
		input.Limit = minSearchProfileLimit
	}

	if input.Offset < 0 {
		input.Offset = 0
	}

	z.Infow("trying with", "input", input)

	output := &dto.SearchProfileOutputDto{
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	total, err := uc.profileService.Count(ctx, input.Query, input.CurrentUserID, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	profiles, err := uc.profileService.FindProfiles(ctx, input.Query, input.CurrentUserID, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	output.Profiles = profiles
	output.Total = total

	return output, nil
}
