package auth

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sasvyn/backend/internal/modules/sessions"
	"github.com/sasvyn/backend/internal/modules/users"
)

type Service struct {
	userRepository *users.Repository
	sessionService *sessions.Service
}

func NewService(
	userRepository *users.Repository,
	sessionService *sessions.Service,
) *Service {
	return &Service{
		userRepository: userRepository,
		sessionService: sessionService,
	}
}

func (s *Service) SocialLogin(
	ctx context.Context,
	request SocialLoginRequest,
) (users.User, string, string, error) {
	start := time.Now()

	userStart := time.Now()
	user, err := s.userRepository.GetByAppleID(ctx, request.AppleID)
	log.Printf("[SocialLogin] GetByAppleID: %v", time.Since(userStart))
	if err != nil {
		if err != sql.ErrNoRows {
			return users.User{}, "", "", err
		}
		createStart := time.Now()
		now := time.Now().UTC()
		user = &users.User{
			ID:        uuid.NewString(),
			AppleID:   request.AppleID,
			FullName:  request.FullName,
			Email:     request.Email,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := s.userRepository.Create(ctx, *user); err != nil {
			return users.User{}, "", "", err
		}
		log.Printf("[SocialLogin] CreateUser: %v", time.Since(createStart))
	}
	sessionStart := time.Now()
	accessToken, refreshToken, err := s.sessionService.CreateSession(ctx, user.ID)
	log.Printf("[SocialLogin] CreateSession: %v", time.Since(sessionStart))
	if err != nil {
		return users.User{}, "", "", err
	}
	log.Printf("[SocialLogin] TOTAL SERVICE: %v", time.Since(start))
	return *user, accessToken, refreshToken, nil
}
