package app

import (
	"database/sql"
	"net/http"

	"github.com/sasvyn/backend/internal/api"
	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/sessions"
	"github.com/sasvyn/backend/internal/modules/skills"
	"github.com/sasvyn/backend/internal/modules/users"
	"github.com/sasvyn/backend/internal/ratelimit"
)

func BuildRouter(db *sql.DB, limiters *RateLimiters) http.Handler {
	sessionRepository := sessions.NewRepository(db)
	sessionService := sessions.NewService(sessionRepository)

	userRepository := users.NewRepository(db)
	skillRepository := skills.NewRepository(db)

	authService := auth.NewService(userRepository, sessionService)
	authHandler := auth.NewHandler(authService, sessionService)
	authMiddleWare := auth.NewMiddleware(sessionService)
	userHandler := users.NewHandler(userRepository)
	skillsHandler := skills.NewHandler(skillRepository)

	router := api.NewRouter(api.Dependencies{
		AuthHandler:    authHandler,
		AuthMiddleware: authMiddleWare,
		UserHandler:    userHandler,
		SkillsHandler:  skillsHandler,
	})

	return ratelimit.PolicyMiddleware(
		limiters.Default,
		limiters.AuthRefresh,
		limiters.SocialLogin,
	)(router)
}
