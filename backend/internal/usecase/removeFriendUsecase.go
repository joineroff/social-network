package usecase

import (
	"context"

	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/service"
	"github.com/joineroff/social-network/backend/internal/validation"
)

type RemoveFriendUsecase struct {
	friendService service.FriendService
}

func NewRemoveFriendUsecase(
	friendService service.FriendService,
) *RemoveFriendUsecase {
	return &RemoveFriendUsecase{
		friendService: friendService,
	}
}

func (uc *RemoveFriendUsecase) Do(ctx context.Context, input *dto.RemoveFriendInputDto) error {
	if err := validation.Struct(input); err != nil {
		return err
	}

	return uc.friendService.RemoveFriend(ctx, input.CurrentUserID, input.FriendID)
}
