package skills

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sasvyn/backend/internal/middleware"
	"github.com/sasvyn/backend/internal/modules/idempotency"
	"github.com/sasvyn/backend/internal/realtime"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	repository         *Repository
	idempotencyService *idempotency.Service
	publisher          *realtime.Publisher
}

func NewHandler(repository *Repository, idempotencyService *idempotency.Service, publisher *realtime.Publisher) *Handler {
	return &Handler{repository: repository, idempotencyService: idempotencyService, publisher: publisher}
}

// Create godoc
// @Summary Create a skill
// @Description Create a new skill for the authenticated user
// @Tags Skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Idempotency-Key header string true "Idempotency key"
// @Param X-Client-ID header string true "Client ID"
// @Param request body CreateSkillDTO true "Skill details"
// @Success 201 {object} map[string]interface{} "Skill created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 409 {object} map[string]interface{} "Conflict"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /skills [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	var request CreateSkillDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			now := time.Now().UTC()

			skill := Skill{
				ID:          request.ID,
				UserID:      userID,
				Skill:       request.Skill,
				Category:    request.Category,
				SyncVersion: 1,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			if err := h.repository.Create(r.Context(), skill); err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusCreated,
				Message:    "Skill added successfully",
				Data:       skill,
			}, nil
		},
	)

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(
				w,
				http.StatusNotFound,
				"Skill was not found",
			)
			return
		}

		log.Printf("Create Skill error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Skill could not be created",
		)
		return
	}

	response.WriteRaw(w, result.StatusCode, result.Body)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventSkillCreated,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish skill.created event: %v", err)
		}
	}
}

// Fetch godoc
// @Summary Fetch all skills
// @Description Fetch all skills belonging to the authenticated user
// @Tags Skills
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Skills fetched successfully"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 500 {object} map[string]interface{} "Failed to fetch skills"
// @Router /skills [get]
func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	skills, err := h.repository.Fetch(r.Context(), userID)
	if err != nil {
		response.Write(w, http.StatusInternalServerError, "failed to fetch skills")
		return
	}

	response.WriteList(w, http.StatusOK, "Skill fetched successfully", skills)
}

// FetchByID godoc
// @Summary Fetch a skill by ID
// @Description Fetch a specific skill belonging to the authenticated user
// @Tags Skills
// @Produce json
// @Security BearerAuth
// @Param id path string true "Skill ID"
// @Success 200 {object} map[string]interface{} "Skill fetched successfully"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 404 {object} map[string]interface{} "Skill not found"
// @Failure 500 {object} map[string]interface{} "Failed to fetch skill"
// @Router /skills/{id} [get]
func (h *Handler) FetchByID(w http.ResponseWriter, r *http.Request) {
	skillID := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())

	skill, err := h.repository.FetchByID(r.Context(), skillID, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "Skill was not found")
			return
		}

		log.Printf("Fetch Skill error: %v", err)
		response.Write(w, http.StatusInternalServerError, "Failed to fetch Skill")
		return
	}

	response.WriteItem(
		w,
		http.StatusOK,
		"Skill fetched successfully",
		skill,
	)
}

// Update godoc
// @Summary Update a skill
// @Description Update an existing skill belonging to the authenticated user
// @Tags Skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Skill ID"
// @Param Idempotency-Key header string true "Idempotency key"
// @Param X-Client-ID header string true "Client ID"
// @Param request body UpdateSkillDTO true "Updated skill details"
// @Success 200 {object} map[string]interface{} "Skill updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 404 {object} map[string]interface{} "Skill not found"
// @Failure 409 {object} map[string]interface{} "Conflict"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /skills/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	skillID := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	var request UpdateSkillDTO
	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		updatedSkill, err := h.repository.Update(r.Context(), skillID, userID, request)
		if err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusOK,
			Message:    "Skill updated successfully",
			Data:       updatedSkill,
		}, nil
	})

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(
				w,
				http.StatusConflict,
				err.Error(),
			)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(
				w,
				http.StatusNotFound,
				"Skill was not found",
			)
			return
		}

		if strings.Contains(err.Error(), "23505") {
			response.Write(
				w,
				http.StatusConflict,
				"Skill already exists",
			)
			return
		}

		log.Printf("Update Skill error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Skill could not be updated",
		)

		return
	}

	response.WriteRaw(
		w,
		result.StatusCode,
		result.Body,
	)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventSkillUpdated,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish skill.created event: %v", err)
		}
	}
}

// Delete godoc
// @Summary Delete a skill
// @Description Delete a skill belonging to the authenticated user
// @Tags Skills
// @Produce json
// @Security BearerAuth
// @Param id path string true "Skill ID"
// @Param Idempotency-Key header string true "Idempotency key"
// @Param X-Client-ID header string true "Client ID"
// @Success 200 {object} map[string]interface{} "Skill deleted successfully"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 404 {object} map[string]interface{} "Skill not found"
// @Failure 409 {object} map[string]interface{} "Conflict"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /skills/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	skillID := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		if err := h.repository.Delete(r.Context(), skillID, userID); err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusOK,
			Message:    "Skill deleted successfully",
		}, nil
	})

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(
				w,
				http.StatusNotFound,
				"Skill was not found",
			)
			return
		}

		log.Printf("Delete Skill error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Skill could not be deleted",
		)
		return
	}

	response.WriteRaw(
		w,
		result.StatusCode,
		result.Body,
	)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventSkillDeleted,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish skill.created event: %v", err)
		}
	}
}
