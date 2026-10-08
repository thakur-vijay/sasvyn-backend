package api

import (
	"net/http"

	"github.com/sasvyn/backend/internal/middleware"
	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/documents"
	"github.com/sasvyn/backend/internal/modules/education"
	"github.com/sasvyn/backend/internal/modules/languages"
	"github.com/sasvyn/backend/internal/modules/skills"
	sociallinks "github.com/sasvyn/backend/internal/modules/socialLinks"
	"github.com/sasvyn/backend/internal/modules/upload"
	"github.com/sasvyn/backend/internal/modules/users"
	"github.com/sasvyn/backend/internal/realtime"
)

const v1Prefix = "/api/v1"

type Dependencies struct {
	AuthHandler        *auth.Handler
	AuthMiddleware     *middleware.Middleware
	ClientIdMiddleware *middleware.ClientIDMiddleware
	UserHandler        *users.Handler
	SkillsHandler      *skills.Handler
	LanguagesHandler   *languages.Handler
	UploadHandler      *upload.Handler
	SocialLinkHandler  *sociallinks.Handler
	RealtimeHandler    *realtime.Handler
	RealtimePublisher  *realtime.Publisher
	DocumentsHandler   *documents.Handler
	EducationsHandler  *education.Handler
}

func NewRouter(dependencies Dependencies) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)

	v1Mux := http.NewServeMux()
	auth.RegisterRoutes(v1Mux, dependencies.AuthHandler, dependencies.AuthMiddleware, dependencies.ClientIdMiddleware)
	users.RegisterRoutes(v1Mux, dependencies.UserHandler, dependencies.AuthMiddleware)
	skills.RegisterRoutes(v1Mux, dependencies.SkillsHandler, dependencies.AuthMiddleware, dependencies.ClientIdMiddleware)
	languages.RegisterRoutes(v1Mux, dependencies.LanguagesHandler, dependencies.AuthMiddleware)
	upload.RegisterRoutes(v1Mux, dependencies.UploadHandler, dependencies.AuthMiddleware)
	sociallinks.RegisterRoutes(v1Mux, dependencies.SocialLinkHandler, dependencies.AuthMiddleware)
	realtime.RegisterRoutes(v1Mux, dependencies.RealtimeHandler, dependencies.AuthMiddleware)
	documents.RegisterRoutes(v1Mux, dependencies.DocumentsHandler, dependencies.AuthMiddleware, dependencies.ClientIdMiddleware)
	education.RegisterRoutes(v1Mux, dependencies.EducationsHandler, dependencies.AuthMiddleware, dependencies.ClientIdMiddleware)
	mux.Handle(v1Prefix+"/", http.StripPrefix(v1Prefix, v1Mux))

	return mux
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
