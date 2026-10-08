package education

import (
	"context"
	"time"
)

type Service struct {
	repository *Repository
}

func NewService(
	repository *Repository,
) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	request CreateEducationDTO,
	userID string,
) (Education, error) {
	now := time.Now().UTC()

	education := Education{
		ID:           request.ID,
		UserID:       userID,
		Degree:       request.Degree,
		FieldOfStudy: request.FieldOfStudy,
		Institution:  request.Institution,
		StartDate:    request.StartDate,
		EndDate:      request.EndDate,
		IsPursuing:   request.IsPursuing,
		Grade:        request.Grade,
		GradeType:    request.GradeType,
		Description:  request.Description,
		SyncVersion:  1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repository.Create(ctx, education); err != nil {
		return Education{}, err
	}

	return education, nil
}

func (s *Service) Fetch(ctx context.Context, userID string) ([]Education, error) {
	educations, err := s.repository.Fetch(ctx, userID)
	if err != nil {
		return nil, err
	}

	return educations, nil
}

func (s *Service) FetchByID(ctx context.Context, educationID, userID string) (Education, error) {
	education, err := s.repository.FetchByID(ctx, educationID, userID)
	if err != nil {
		return Education{}, err
	}

	return education, nil
}

func (s *Service) Update(ctx context.Context, request UpdateEducationDTO, educationID, userID string) (Education, error) {
	return s.repository.Update(ctx, educationID, userID, request)
}

func (s *Service) Delete(
	ctx context.Context,
	educationID string,
	userID string,
) error {
	return s.repository.Delete(ctx, educationID, userID)
}
