package repository

import (
	"context"

	"github.com/joineroff/social-network/backend/internal/entity"
)

type ProfileRepository interface {
	// count profiles satisfied to searchQuery
	GetProfile(
		ctx context.Context,
		id string,
		currentUserID string,
	) (*entity.Profile, error)
	// count profiles satisfied to searchQuery
	CountProfiles(
		ctx context.Context,
		searchQuery string,
		currentUserID string,
		requreFriend bool,
	) (int, error)
	// find profiles satisfied to searchQuery
	FindProfiles(
		ctx context.Context,
		searchQuery string,
		currentUserID string,
		requreFriend bool,
		limit int,
		offset int,
	) ([]*entity.Profile, error)
}
