package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/entity"
	"github.com/joineroff/social-network/backend/internal/service"
	"github.com/joineroff/social-network/backend/internal/validation"
	"go.uber.org/zap"
)

type SignUpUsecase struct {
	authService service.AuthService
	userService service.UserService
}

func NewSignUpUsecase(
	authService service.AuthService,
	userService service.UserService,
) *SignUpUsecase {
	return &SignUpUsecase{
		authService: authService,
		userService: userService,
	}
}

func (uc *SignUpUsecase) Do(ctx context.Context, input *dto.SignUpInputDto) (*dto.SignUpOutputDto, error) {
	z := zap.S().With("context", "SignUpUsecase")

	output := &dto.SignUpOutputDto{}

	if err := validation.Struct(input); err != nil {
		return nil, err
	}

	user, err := uc.userService.FindByLogin(ctx, input.Login)
	if err != nil && !errors.Is(err, service.ErrUserNotFound) {
		z.Error(err)
		return nil, err
	}

	if user != nil {
		z.Error(err)
		return nil, NewValidationError("provided login name is busy")
	}

	user = &entity.User{
		Login:     input.Login,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		City:      input.City,
		Interests: input.Interests,
		Gender:    input.Gender,
	}

	passHash, err := uc.authService.CreatePasswordHash(input.Password)
	if err != nil {
		z.Error(err)
		return nil, fmt.Errorf("something went wrong")
	}

	user.Password = passHash

	user, err = uc.userService.Create(ctx, user)
	if err != nil {
		z.Error(err)
		return nil, err
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
