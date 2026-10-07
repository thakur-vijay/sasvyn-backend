package documents

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sasvyn/backend/internal/storage"
)

type Service struct {
	repository *Repository
	r2Client   *s3.Client
}

func NewService(
	repository *Repository,
	r2Client *s3.Client,
) *Service {
	return &Service{
		repository: repository,
		r2Client:   r2Client,
	}
}

func (s *Service) Create(
	ctx context.Context,
	request CreateDocumentDTO,
	userID string,
) (Document, error) {
	now := time.Now().UTC()

	document := Document{
		ID:          request.ID,
		UserID:      userID,
		Name:        request.Name,
		Category:    request.Category,
		FileSize:    request.FileSize,
		Key:         &request.Key,
		SyncVersion: 1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repository.Create(ctx, document); err != nil {
		return Document{}, err
	}

	return document, nil
}

func (s *Service) Fetch(
	ctx context.Context,
	userID string,
) ([]Document, error) {
	return s.repository.Fetch(ctx, userID)
}

func (s *Service) FetchByID(
	ctx context.Context,
	documentID string,
	userID string,
) (Document, error) {
	return s.repository.FetchByID(ctx, documentID, userID)
}

func (s *Service) Delete(
	ctx context.Context,
	documentID string,
	userID string,
) error {
	document, err := s.repository.FetchByID(ctx, documentID, userID)
	if err != nil {
		return err
	}

	if document.Key != nil {
		if err := storage.DeleteObject(ctx, s.r2Client, *document.Key); err != nil {
			return err
		}
	}

	return s.repository.Delete(ctx, documentID, userID)
}
