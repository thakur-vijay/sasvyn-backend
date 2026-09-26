package app

import (
	"database/sql"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sasvyn/backend/internal/api"
	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/languages"
	"github.com/sasvyn/backend/internal/modules/sessions"
	"github.com/sasvyn/backend/internal/modules/skills"
	sociallinks "github.com/sasvyn/backend/internal/modules/socialLinks"
	"github.com/sasvyn/backend/internal/modules/upload"
	"github.com/sasvyn/backend/internal/modules/users"
	"github.com/sasvyn/backend/internal/ratelimit"
)

func BuildRouter(db *sql.DB, limiters *RateLimiters, r2Client *s3.Client) http.Handler {
	presignClient := s3.NewPresignClient(r2Client)
	sessionRepository := sessions.NewRepository(db)
	sessionService := sessions.NewService(sessionRepository)

	userRepository := users.NewRepository(db)
	skillRepository := skills.NewRepository(db)
	languagesRepository := languages.NewRepository(db)
	socialLinksRepository := sociallinks.NewRepository(db)

	authService := auth.NewService(userRepository, sessionService)
	authHandler := auth.NewHandler(authService, sessionService)
	authMiddleWare := auth.NewMiddleware(sessionService)
	userHandler := users.NewHandler(userRepository, presignClient)
	skillsHandler := skills.NewHandler(skillRepository)
	languagesHandler := languages.NewHandler(languagesRepository)
	uploadHandler := upload.NewHandler(r2Client, presignClient)
	socialLinksHandler := sociallinks.NewHandler(socialLinksRepository)

	router := api.NewRouter(api.Dependencies{
		AuthHandler:       authHandler,
		AuthMiddleware:    authMiddleWare,
		UserHandler:       userHandler,
		SkillsHandler:     skillsHandler,
		LanguagesHandler:  languagesHandler,
		UploadHandler:     uploadHandler,
		SocialLinkHandler: socialLinksHandler,
	})

	return ratelimit.PolicyMiddleware(
		limiters.Default,
		limiters.AuthRefresh,
		limiters.SocialLogin,
	)(router)
}
