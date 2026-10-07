package app

import (
	"database/sql"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sasvyn/backend/internal/api"
	"github.com/sasvyn/backend/internal/middleware"
	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/idempotency"
	"github.com/sasvyn/backend/internal/modules/languages"
	"github.com/sasvyn/backend/internal/modules/sessions"
	"github.com/sasvyn/backend/internal/modules/skills"
	sociallinks "github.com/sasvyn/backend/internal/modules/socialLinks"
	"github.com/sasvyn/backend/internal/modules/upload"
	"github.com/sasvyn/backend/internal/modules/users"
	"github.com/sasvyn/backend/internal/ratelimit"
	"github.com/sasvyn/backend/internal/realtime"
)

func BuildRouter(db *sql.DB, limiters *RateLimiters, r2Client *s3.Client) http.Handler {
	// Storage
	presignClient := s3.NewPresignClient(r2Client)

	// Sessions
	sessionRepository := sessions.NewRepository(db)
	sessionService := sessions.NewService(sessionRepository)

	// Idempotency
	idempotencyRepository := idempotency.NewRepository(db)
	idempotencyService := idempotency.NewService(idempotencyRepository)

	// Repositories
	userRepository := users.NewRepository(db)
	skillRepository := skills.NewRepository(db)
	languagesRepository := languages.NewRepository(db)
	socialLinksRepository := sociallinks.NewRepository(db)

	// Realtime
	realtimeManager := realtime.NewManager()
	realtimeHandler := realtime.NewHandler(realtimeManager)
	realtimePublisher := realtime.NewPublisher(realtimeManager)

	// Auth
	authService := auth.NewService(userRepository, sessionService)
	authHandler := auth.NewHandler(authService, sessionService)

	//Middlewares
	authMiddleware := middleware.NewAuthMiddleware(sessionService)
	clientIdMiddleware := middleware.NewClientIDMiddleware()

	// Handlers
	userHandler := users.NewHandler(userRepository, presignClient, idempotencyService)
	skillsHandler := skills.NewHandler(skillRepository, idempotencyService, realtimePublisher)
	languagesHandler := languages.NewHandler(languagesRepository, idempotencyService)
	uploadHandler := upload.NewHandler(r2Client, presignClient)
	socialLinksHandler := sociallinks.NewHandler(socialLinksRepository, idempotencyService)

	// Router
	router := api.NewRouter(api.Dependencies{
		AuthHandler:        authHandler,
		AuthMiddleware:     authMiddleware,
		ClientIdMiddleware: clientIdMiddleware,
		UserHandler:        userHandler,
		SkillsHandler:      skillsHandler,
		LanguagesHandler:   languagesHandler,
		UploadHandler:      uploadHandler,
		SocialLinkHandler:  socialLinksHandler,
		RealtimeHandler:    realtimeHandler,
		RealtimePublisher:  realtimePublisher,
	})

	// Global Middleware
	return ratelimit.PolicyMiddleware(
		limiters.Default,
		limiters.AuthRefresh,
		limiters.SocialLogin,
	)(router)
}
