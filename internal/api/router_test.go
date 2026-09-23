package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/sessions"
	"github.com/sasvyn/backend/internal/modules/users"
)

func TestNewRouter_ComposesHealthAndVersionedFeatureRoutes(t *testing.T) {
	userRepository := users.NewRepository(nil)
	sessionService := sessions.NewService(sessions.NewRepository(nil))
	authMiddleware := auth.NewMiddleware(sessionService)

	router := NewRouter(Dependencies{
		AuthHandler:    auth.NewHandler(auth.NewService(userRepository, sessionService), sessionService),
		AuthMiddleware: authMiddleware,
		UserHandler:    users.NewHandler(userRepository, nil),
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "health is not versioned",
			method:         http.MethodGet,
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "versioned user route is protected",
			method:         http.MethodGet,
			path:           "/api/v1/users/user-id",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "unversioned user route is not registered",
			method:         http.MethodGet,
			path:           "/users/user-id",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "health is not mounted under the API",
			method:         http.MethodGet,
			path:           "/api/v1/health",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != test.expectedStatus {
				t.Fatalf("expected status %d, got %d", test.expectedStatus, recorder.Code)
			}
		})
	}
}
