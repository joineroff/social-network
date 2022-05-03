package service

import (
	"context"

	"github.com/joineroff/social-network/backend/internal/repository"
)

var _ FriendService = &friendService{}

type FriendService interface {
	CountFriends(ctx context.Context, userID string) (int, error)
	AddFriend(ctx context.Context, userID string, friendID string) error
	RemoveFriend(ctx context.Context, userID string, friendID string) error
}

type friendService struct {
	friendRepository repository.FriendRepository
}

func NewFriendService(
	friendRepository repository.FriendRepository,
) *friendService {
	return &friendService{
		friendRepository: friendRepository,
	}
}

func (s *friendService) CountFriends(ctx context.Context, userID string) (int, error) {
	return s.friendRepository.Count(ctx, userID)
}

func (s *friendService) AddFriend(ctx context.Context, userID, friendID string) error {
	return s.friendRepository.Add(ctx, userID, friendID)
}

func (s *friendService) RemoveFriend(ctx context.Context, userID, friendID string) error {
	return s.friendRepository.Remove(ctx, userID, friendID)
}
