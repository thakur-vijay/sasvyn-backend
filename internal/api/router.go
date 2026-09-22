package api

import (
	"net/http"

	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/languages"
	"github.com/sasvyn/backend/internal/modules/skills"
	"github.com/sasvyn/backend/internal/modules/upload"
	"github.com/sasvyn/backend/internal/modules/users"
)

const v1Prefix = "/api/v1"

type Dependencies struct {
	AuthHandler      *auth.Handler
	AuthMiddleware   *auth.Middleware
	UserHandler      *users.Handler
	SkillsHandler    *skills.Handler
	LanguagesHandler *languages.Handler
	UploadHandler    *upload.Handler
}

func NewRouter(dependencies Dependencies) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)

	v1Mux := http.NewServeMux()
	auth.RegisterRoutes(v1Mux, dependencies.AuthHandler, dependencies.AuthMiddleware)
	users.RegisterRoutes(v1Mux, dependencies.UserHandler, dependencies.AuthMiddleware.RequireAuth)
	skills.RegisterRoutes(v1Mux, dependencies.SkillsHandler, dependencies.AuthMiddleware.RequireAuth)
	languages.RegisterRoutes(v1Mux, dependencies.LanguagesHandler, dependencies.AuthMiddleware.RequireAuth)
	upload.RegisterRoutes(v1Mux, dependencies.UploadHandler, dependencies.AuthMiddleware.RequireAuth)

	mux.Handle(v1Prefix+"/", http.StripPrefix(v1Prefix, v1Mux))

	return mux
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
