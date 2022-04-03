package usecase

import (
	"context"
	"errors"

	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/service"
	"github.com/joineroff/social-network/backend/internal/validation"
	"go.uber.org/zap"
)

type SignInUsecase struct {
	authService service.AuthService
	userService service.UserService
}

func NewSignInUsecase(
	authService service.AuthService,
	userService service.UserService,
) *SignInUsecase {
	return &SignInUsecase{
		authService: authService,
		userService: userService,
	}
}

func (uc *SignInUsecase) Do(ctx context.Context, input *dto.SignInInputDto) (*dto.SignInOutputDto, error) {
	z := zap.S().With("context", "SignUpUsecase")
	output := &dto.SignInOutputDto{}

	if err := validation.Struct(input); err != nil {
		return nil, err
	}

	user, err := uc.userService.FindByLogin(ctx, input.Login)
	if err != nil && !errors.Is(err, service.ErrUserNotFound) {
		z.Error(err)
		return nil, err
	} else if errors.Is(err, service.ErrUserNotFound) {
		return nil, NewValidationError("invalid credentials")
	}

	err = uc.authService.ComparePasswordWithHash(input.Password, user.Password)
	if err != nil {
		z.Error(err)
		return nil, NewValidationError("invalid credentials")
	}

	data, err := uc.authService.GenerateTokens(user)
	if err != nil {
		z.Error(err)
		return nil, err
	}

	output.AccessToken = data.Access
	output.RefreshToken = data.Refresh

	return output, nil
}
