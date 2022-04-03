package service

import (
	"context"

	"github.com/joineroff/social-network/backend/internal/entity"
	"github.com/joineroff/social-network/backend/internal/repository"
)

type ProfileService interface {
	// FindProfiles ...
	GetProfile(
		ctx context.Context,
		userID string,
		currentUserID string,
	) (*entity.Profile, error)
	// FindProfiles ...
	FindProfiles(
		ctx context.Context,
		searchQuery string,
		currentUserID string,
		limit int,
		offset int,
	) ([]*entity.Profile, error)

	// FindFriends ...
	FindFriends(
		ctx context.Context,
		searchQuery string,
		currentUserID string,
		limit int,
		offset int,
	) ([]*entity.Profile, error)

	// Count ...
	Count(
		ctx context.Context,
		searchQuery string,
		currentUserID string,
		limit int,
		offset int,
	) (int, error)

	// CountFriends ...
	CountFriends(
		ctx context.Context,
		searchQuery string,
		currentUserID string,
	) (int, error)
}

type profileService struct {
	profileRepository repository.ProfileRepository
}

func NewProfileService(
	profileRepository repository.ProfileRepository,
) ProfileService {
	return &profileService{
		profileRepository: profileRepository,
	}
}

func (s *profileService) GetProfile(
	ctx context.Context,
	userID string,
	currentUserID string,
) (*entity.Profile, error) {
	return s.profileRepository.GetProfile(ctx, userID, currentUserID)
}

func (s *profileService) FindProfiles(
	ctx context.Context,
	searchQuery string,
	currentUserID string,
	limit int,
	offset int,
) ([]*entity.Profile, error) {
	return s.profileRepository.FindProfiles(ctx, searchQuery, currentUserID, false, limit, offset)
}

func (s *profileService) FindFriends(
	ctx context.Context,
	searchQuery string,
	currentUserID string,
	limit int,
	offset int,
) ([]*entity.Profile, error) {
	return s.profileRepository.FindProfiles(ctx, searchQuery, currentUserID, true, limit, offset)
}

func (s *profileService) Count(
	ctx context.Context,
	searchQuery string,
	currentUserID string,
	limit int,
	offset int,
) (int, error) {
	return s.profileRepository.CountProfiles(ctx, searchQuery, currentUserID, false)
}

func (s *profileService) CountFriends(
	ctx context.Context,
	searchQuery string,
	currentUserID string,
) (int, error) {
	return s.profileRepository.CountProfiles(ctx, searchQuery, currentUserID, true)
}
