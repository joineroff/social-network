package service

import "context"

var _ SubscriptionService = &subscriptionService{}

type SubscriptionService interface {
	Subscribe(ctx context.Context, userID string, toUser string) error
	Unsubscribe(ctx context.Context, userID string, fromUser string) error
	CountSubscribers(ctx context.Context, userID string) (int, error)
	CountSubscriptions(ctx context.Context, userID string) (int, error)
	GetSubscribers(ctx context.Context, userID string) ([]string, error)
	GetSubscriptions(ctx context.Context, userID string) ([]string, error)
}

type subscriptionService struct{}

func NewSubscriptionService() *subscriptionService {
	return &subscriptionService{}
}

func (s *subscriptionService) Subscribe(ctx context.Context, userID, toUser string) error {
	panic("not implemented") // TODO: Implement
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, userID, fromUser string) error {
	panic("not implemented") // TODO: Implement
}

func (s *subscriptionService) CountSubscribers(ctx context.Context, userID string) (int, error) {
	panic("not implemented") // TODO: Implement
}

func (s *subscriptionService) CountSubscriptions(ctx context.Context, userID string) (int, error) {
	panic("not implemented") // TODO: Implement
}

func (s *subscriptionService) GetSubscribers(ctx context.Context, userID string) ([]string, error) {
	panic("not implemented") // TODO: Implement
}

func (s *subscriptionService) GetSubscriptions(ctx context.Context, userID string) ([]string, error) {
	panic("not implemented") // TODO: Implement
}
