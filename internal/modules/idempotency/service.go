package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"uuid"

	"github.com/sasvyn/backend/internal/response"
)

type Result struct {
	StatusCode int
	Message    string
	Data       any
	Body       []byte
	Replayed   bool
}

var ErrAlreadyProcessing = errors.New(

	"request with this Idempotency-Key is already being processed",
)

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
) (*IdempotencyRecord, bool, error) {
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

func (s *Service) Execute(
	ctx context.Context,
	userID string,
	key string,
	operation func() (Result, error),
) (Result, error) {
	record := IdempotencyRecord{
		ID:             uuid.New().String(),
		UserID:         userID,
		IdempotencyKey: key,
		CreatedAt:      time.Now().UTC(),
	}

	reserved, ownsReservation, err := s.Reserve(ctx, record)
	if err != nil {
		return Result{}, err
	}

	if reserved.Status == "completed" {
		return Result{
			StatusCode: *reserved.StatusCode,
			Body:       reserved.ResponseBody,
		}, nil
	}

	if !ownsReservation {
		return Result{}, ErrAlreadyProcessing
	}

	if reserved.Status == "processing" {
		return Result{}, ErrAlreadyProcessing
	}

	time.Sleep(3 * time.Second)
	result, err := operation()
	if err != nil {
		_ = s.Delete(ctx, userID, key)
		return Result{}, err
	}

	responseBody, err := json.Marshal(
		response.Response{
			StatusCode: result.StatusCode,
			Message:    result.Message,
			Data:       result.Data,
		},
	)
	if err != nil {
		_ = s.Delete(ctx, userID, key)
		return Result{}, err
	}

	if err := s.Complete(
		ctx,
		userID,
		key,
		result.StatusCode,
		responseBody,
	); err != nil {
		return Result{}, err
	}

	result.Body = responseBody

	return result, nil
}
