package repository

import (
	"context"

	"github.com/joineroff/social-network/backend/internal/entity"
)

type FindUserCriteria struct {
	ID    string
	Login string
	Name  string
	City  string
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, model *entity.User) (*entity.User, error)
	Update(ctx context.Context, model *entity.User) (*entity.User, error)
	Delete(ctx context.Context, model *entity.User) error
}
