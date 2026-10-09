package experiences

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
	request CreateExperienceDTO,
	userID string,
) (Experience, error) {
	now := time.Now().UTC()

	experience := Experience{
		ID:                 request.ID,
		UserID:             userID,
		Role:               request.Role,
		Company:            request.Company,
		Location:           request.Location,
		StartDate:          request.StartDate,
		EndDate:            request.EndDate,
		IsCurrentlyWorking: request.IsCurrentlyWorking,
		SyncVersion:        1,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repository.Create(ctx, experience); err != nil {
		return Experience{}, err
	}

	return experience, nil
}

func (s *Service) Fetch(ctx context.Context, userID string) ([]Experience, error) {
	experiences, err := s.repository.Fetch(ctx, userID)
	if err != nil {
		return nil, err
	}

	return experiences, nil
}

func (s *Service) FetchByID(ctx context.Context, id, userID string) (Experience, error) {
	experience, err := s.repository.FetchByID(ctx, id, userID)
	if err != nil {
		return Experience{}, err
	}

	return experience, nil
}

func (s *Service) Update(ctx context.Context, request UpdateExperienceDTO, id, userID string) (Experience, error) {
	return s.repository.Update(ctx, id, userID, request)
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {
	return s.repository.Delete(ctx, id, userID)
}

func (s *Service) CreateResponsibility(
	ctx context.Context,
	request CreateExperienceResponsibilityDTO,
	experienceID string,
	userID string,
) (ExperienceResponsibility, error) {
	now := time.Now().UTC()

	responsibility := ExperienceResponsibility{
		ID:             request.ID,
		ExperienceID:   experienceID,
		Responsibility: request.Responsibility,
		Order:          request.Order,
		SyncVersion:    1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repository.CreateResponsibility(ctx, responsibility, userID); err != nil {
		return ExperienceResponsibility{}, err
	}

	return responsibility, nil
}

func (s *Service) FetchResponsibilities(
	ctx context.Context,
	experienceID string,
	userID string,
) ([]ExperienceResponsibility, error) {
	experience, err := s.repository.FetchByID(ctx, experienceID, userID)
	if err != nil {
		return nil, err
	}

	return s.repository.FetchResponsibilities(ctx, experience.ID)
}

func (s *Service) UpdateResponsibility(
	ctx context.Context,
	request UpdateExperienceResponsibilityDTO,
	userID, experienceID, responsibilityID string,
) (ExperienceResponsibility, error) {
	return s.repository.UpdateResponsibility(
		ctx,
		userID,
		experienceID,
		responsibilityID,
		request,
	)
}

func (s *Service) DeleteResponsibility(
	ctx context.Context,
	userID, experienceID, responsibilityID string,
) error {
	return s.repository.DeleteResponsibility(
		ctx,
		userID,
		experienceID,
		responsibilityID,
	)
}
