package idempotency

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Find(
	ctx context.Context,
	userID string,
	key string,
) (*IdempotencyRecord, error) {
	return s.repository.Find(ctx, userID, key)
}

func (s *Service) Save(
	ctx context.Context,
	record IdempotencyRecord,
) error {
	return s.repository.Save(ctx, record)
}

func (s *Service) Reserve(
	ctx context.Context,
	record IdempotencyRecord,
) (*IdempotencyRecord, error) {
	return s.repository.Reserve(ctx, record)
}

func (s *Service) Complete(
	ctx context.Context,
	userID string,
	key string,
	statusCode int,
	responseBody []byte,
) error {
	return s.repository.Complete(
		ctx,
		userID,
		key,
		statusCode,
		responseBody,
	)
}

func (s *Service) Delete(
	ctx context.Context,
	userID string,
	key string,
) error {
	return s.repository.Delete(ctx, userID, key)
}
