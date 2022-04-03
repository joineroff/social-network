package repository

import (
	"context"
)

type FriendRepository interface {
	// Add friendID as friend to userID
	// If already exists ignore
	Add(ctx context.Context, userID, friendID string) error

	// Remove user's friend
	// If not exists ignore
	Remove(ctx context.Context, userID, friendID string) error

	// Count user's friends
	Count(ctx context.Context, userID string) (int, error)
}
