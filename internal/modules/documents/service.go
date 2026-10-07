package documents

import (
	"context"

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
	document Document,
) error {
	return s.repository.Create(ctx, document)
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
