package usecase

import (
	"context"

	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/service"
	"github.com/joineroff/social-network/backend/internal/validation"
)

type AddFriendUsecase struct {
	friendService service.FriendService
}

func NewAddFriendUsecase(
	friendService service.FriendService,
) *AddFriendUsecase {
	return &AddFriendUsecase{
		friendService: friendService,
	}
}

func (uc *AddFriendUsecase) Do(ctx context.Context, input *dto.AddFriendInputDto) error {
	if err := validation.Struct(input); err != nil {
		return err
	}

	return uc.friendService.AddFriend(ctx, input.CurrentUserID, input.FriendID)
}
