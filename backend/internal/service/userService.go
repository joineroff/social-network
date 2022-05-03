package service

import (
	"context"
	"errors"

	"github.com/joineroff/social-network/backend/internal/entity"
	"github.com/joineroff/social-network/backend/internal/repository"
)

var _ UserService = &userService{}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exists")
)

type FindUserCriteria struct {
	ID    string
	Login string
	Name  string
	City  string
}

type UserService interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(
	userRepository repository.UserRepository,
) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u *userService) FindByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	return nil, ErrUserNotFound
}

func (u *userService) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	user, err := u.userRepository.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	return nil, ErrUserNotFound
}

func (u *userService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	foundUser, err := u.FindByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	if foundUser != nil {
		return nil, ErrUserAlreadyExist
	}

	return u.userRepository.Create(ctx, user)
}

func (u *userService) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	foundUser, err := u.FindByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// update all fields except email and password
	// @TODO implement separate UpdateLogin and UpdatePassword logic
	user.Password = foundUser.Password
	user.Login = foundUser.Login

	return u.userRepository.Update(ctx, user)
}

func (u *userService) Delete(ctx context.Context, user *entity.User) error {
	_, err := u.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}

	return u.userRepository.Delete(ctx, user)
}
